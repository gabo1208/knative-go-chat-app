package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"

	obsclient "github.com/cloudevents/sdk-go/observability/opencensus/v2/client"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	ceclient "github.com/cloudevents/sdk-go/v2/client"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"knative.dev/eventing/pkg/kncloudevents"
)

var (
	manager  clientManager
	ceClient ceclient.Client
)

type clientManager struct {
	clients    map[*client]string
	usernames  map[string]*client
	broadcast  chan interface{}
	register   chan *client
	unregister chan *client
}

func (manager *clientManager) start() {
	// initialize client to send CloudEvents via HTTP
	var err error
	ceClient, err = obsclient.NewClientHTTP(
		[]cehttp.Option{cehttp.WithIsRetriableFunc(
			func(statusCode int) bool {
				retry, _ := kncloudevents.SelectiveRetry(
					context.TODO(),
					&http.Response{StatusCode: statusCode},
					nil,
				)
				return retry
			})}, nil)
	if err != nil {
		log.Printf("Error creating ceclient %s", err)
		return
	}

	// Start waiting for messages from the different channels
	for {
		select {
		// client connected
		case conn := <-manager.register:
			manager.clients[conn] = conn.id
			manager.usernames[conn.id] = conn
		// client disconnected
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.usernames, conn.id)
				delete(manager.clients, conn)
			}
		// msg to be broadcasted locally in this chat app instance
		case message := <-manager.broadcast:
			log.Printf("Broadasting to %d clients: %+v", len(manager.usernames), message)
			for _, client := range manager.usernames {
				if client.registered {
					log.Printf("Broadasting message to client %s", client.id)
					client.SendCE(message.(*cloudevents.Event), false)
				}
			}
		}
	}
}

func (manager *clientManager) send(message interface{}) {
	for conn := range manager.clients {
		conn.send <- message
	}
}

// Client type for the chat
type client struct {
	id, username string
	registered   bool
	socket       *websocket.Conn
	send         chan interface{}
}

// Handle a new websocket connection
func (c *client) handleWSConnection() {
	defer func() {
		manager.unregister <- c
		if c.registered {
			c.SendCE(
				c.createCE(
					UserDisconnected,
					cloudevents.TextPlain,
					c.username,
				),
				true)
		}
		c.registered = false
		c.socket.Close()
	}()

	// listen to messages in the new socket
	buff := make([]byte, 512)
	for {
		n, err := c.socket.Read(buff)
		if err != nil {
			return
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(buff[:n], &msg); err != nil {
			log.Print(err)
		}
		log.Printf("receiving %s", msg)

		// Check msg type TODO: Going to refactor to ce format
		if val, ok := msg["username"]; ok {
			username := val.(string)
			if reconnecting, ok := msg["reconnecting"]; ok && reconnecting.(bool) {
				log.Printf("reconnecting %s", username)
				if oldClient, ok := manager.usernames[username]; !ok || oldClient == nil {
					log.Println("failed to reconnect")
					return
				}

				c.SendCE(
					c.createCE(
						UserReconnected,
						cloudevents.ApplicationJSON,
						map[string]interface{}{
							"connectedUsers": GetUsernames(manager.usernames),
						},
					),
					false)
				c.connectToClient(username)
			} else if !c.registered {
				log.Printf("First message from client %s updating username", c.id)
				_, exists := manager.usernames[username]
				if exists {
					websocket.Message.Send(c.socket, "error: username already exists")
					continue
				}

				_, ok := manager.clients[c]
				if !ok {
					websocket.Message.Send(c.socket, "error: user does not have a valid connection")
					return
				}

				c.SendCE(
					c.createCE(
						FirstUserConnection,
						cloudevents.ApplicationJSON,
						map[string]interface{}{
							"username":       username,
							"connectedUsers": GetUsernames(manager.usernames),
						},
					),
					false)
				c.connectToClient(username)
			}
		} else {
			if val, ok := manager.usernames[msg["to"].(string)]; ok && val != nil {
				websocket.JSON.Send(
					val.socket,
					map[string]string{
						"from":    c.username,
						"message": msg["message"].(string),
					},
				)
			}
		}
	}
}

// Initialize client manager, init runs only once when the
// controller package is initialized
func init() {
	manager = clientManager{
		broadcast:  make(chan interface{}, 100),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]string),
		usernames:  make(map[string]*client),
	}
	go manager.start()
}

func (c *Controller) WSHandler(ws *websocket.Conn) {
	log.Println("WS connection...")
	client := &client{
		id:     uuid.New().String(),
		socket: ws,
		send:   make(chan interface{}),
	}

	manager.register <- client
	client.handleWSConnection()
}

func GetUsernames(usernamesMap map[string]*client) []string {
	keys := make([]string, 0, len(usernamesMap))
	for k, client := range usernamesMap {
		if client != nil && client.registered {
			keys = append(keys, k)
		}
	}

	return keys
}

// Clients without username get its connection id assigned, here
// this is getting cleaned for the selected unique username
func (c *client) connectToClient(username string) {
	delete(manager.usernames, c.id)
	manager.usernames[username] = c
	c.username = username
	c.registered = true
	c.SendCE(
		c.createCE(
			NewUserConnected,
			cloudevents.TextPlain,
			username,
		),
		true)
}

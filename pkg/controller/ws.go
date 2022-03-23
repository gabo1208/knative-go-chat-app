package controller

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	ceclient "github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/event"
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
			log.Printf("Broadcasting to %d clients: %+v", len(manager.usernames), message)
			for _, client := range manager.usernames {
				if client.registered && client.socket != nil {
					log.Printf("Broadcasting message to client %s", client.id)
					client.SendCE(message.(*cloudevents.Event), false, false)
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
				createCE(
					UserDisconnected,
					cloudevents.TextPlain,
					c.username,
				),
				true,
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
			continue
		}
		log.Printf("receiving %s", msg)
		err = c.processWSMessage(msg, true)
		if err != nil {
			log.Print(err)
			return
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

func (m *clientManager) getExternalUsernames() {
	ce := createCE(GetUsers, cloudevents.TextPlain, os.Getenv("OWN_BROKER_URI"))
	otherClusters := os.Getenv("CLUSTERS_BROKERS_URI")
	for _, uri := range strings.Split(otherClusters, ",") {
		log.Printf("sending GetUsers event to %s", uri)
		SendCEViaHTTP(ce, uri)
	}
}

func AppendUsernames(ce *event.Event, manager *clientManager) {
	var usernames []string
	if err := json.Unmarshal(ce.Data(), &usernames); err != nil {
		log.Print(err)
		return
	}

	for _, username := range usernames {
		manager.usernames[username] = &client{
			id:         uuid.New().String(),
			username:   username,
			registered: true,
		}
	}
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
		createCE(
			NewUserConnected,
			cloudevents.TextPlain,
			username,
		),
		true,
		true)
}

func (c *client) processWSMessage(msg map[string]interface{}, local bool) error {
	// Check msg type TODO: Refactor ws comunication to CE format
	if val, ok := msg["username"]; ok {
		username := val.(string)
		if reconnecting, ok := msg["reconnecting"]; ok && reconnecting.(bool) {
			log.Printf("reconnecting %s", username)
			if oldClient, ok := manager.usernames[username]; !ok || oldClient == nil {
				return errors.New("failed to reconnect")
			}

			c.SendCE(
				createCE(
					UserReconnected,
					cloudevents.ApplicationJSON,
					map[string]interface{}{
						"connectedUsers": GetUsernames(manager.usernames),
					},
				),
				false,
				false)
			c.connectToClient(username)
		} else if !c.registered {
			log.Printf("First message from client %s updating username", c.id)
			_, exists := manager.usernames[username]
			if exists {
				websocket.Message.Send(c.socket, "error: username already exists")
				return nil
			}

			_, ok := manager.clients[c]
			if !ok {
				websocket.Message.Send(c.socket, "error: user does not have a valid connection")
				return errors.New("error: user does not have a valid connection")
			}

			c.SendCE(
				createCE(
					FirstUserConnection,
					cloudevents.ApplicationJSON,
					map[string]interface{}{
						"username":       username,
						"connectedUsers": GetUsernames(manager.usernames),
					},
				),
				false,
				false)
			c.connectToClient(username)
		}
	} else {
		if _, ok := msg["from"]; !ok {
			msg["from"] = c.username
		}

		val, ok := manager.usernames[msg["to"].(string)]
		if ok && val != nil && val.socket != nil {
			websocket.JSON.Send(
				val.socket,
				map[string]string{
					"from":    msg["from"].(string),
					"message": msg["message"].(string),
				},
			)
		} else if (!ok || val.socket == nil) && local {
			c.SendCE(
				createCE(
					MessageFromUser,
					cloudevents.ApplicationJSON,
					msg,
				),
				true,
				false)
		}
	}

	return nil
}

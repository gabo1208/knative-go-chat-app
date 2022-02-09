package controller

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

var (
	manager clientManager
)

type clientManager struct {
	clients    map[*client]string
	usernames  map[string]*client
	broadcast  chan interface{}
	register   chan *client
	unregister chan *client
}

func (manager *clientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = conn.id
			manager.usernames[conn.id] = conn
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.usernames, conn.id)
				delete(manager.clients, conn)
			}
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

type client struct {
	id, username string
	registered   bool
	socket       *websocket.Conn
	send         chan interface{}
}

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

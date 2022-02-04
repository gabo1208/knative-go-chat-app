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
				delete(manager.usernames, manager.clients[conn])
				delete(manager.clients, conn)
			}
		case message := <-manager.broadcast:
			log.Printf("Broadasting to %d clients: %+v", len(manager.clients), message)
			for conn := range manager.clients {
				log.Printf("Broadasting message to client %s", conn.id)
				select {
				case conn.send <- message:
				default:
					log.Print("closing")
					close(conn.send)
					delete(manager.usernames, manager.clients[conn])
					delete(manager.clients, conn)
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
	socket       *websocket.Conn
	send         chan interface{}
}

func (c *client) write() {
	defer func() {
		manager.unregister <- c
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

		if username, ok := msg["username"]; ok {
			usr := username.(string)
			log.Printf("First message from client %s updating username", c.id)
			_, exists := manager.usernames[usr]
			if exists {
				websocket.Message.Send(c.socket, "error: username already exists")
			}

			val, ok := manager.clients[c]
			if !ok {
				websocket.Message.Send(c.socket, "error: user does not have a valid connection")
				return
			}

			delete(manager.usernames, val)
			manager.usernames[usr] = c
			c.username = usr
			//c.SendCE(NewUserConnected, cloudevents.TextPlain, message, true)
			c.SendCE(
				NewUserConnected,
				cloudevents.ApplicationJSON,
				map[string]interface{}{
					"username":       usr,
					"connectedUsers": GetUsernames(manager.usernames),
				},
				false,
			)
		} else {
			websocket.JSON.Send(
				manager.usernames[msg["to"].(string)].socket,
				map[string]string{
					"from":    c.username,
					"message": msg["message"].(string),
				},
			)
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
	client.write()
}

func GetUsernames(usernamesMap map[string]*client) []string {
	keys := make([]string, 0, len(usernamesMap))
	for k := range usernamesMap {
		keys = append(keys, k)
	}

	return keys
}

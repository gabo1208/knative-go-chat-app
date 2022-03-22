package controller

import (
	"log"
	"net/http"
	"sync"

	ceclient "github.com/cloudevents/sdk-go/v2/client"
	"golang.org/x/net/websocket"
)

type Controller struct {
	rootHandler http.Handler
	root        string
	mux         *http.ServeMux
	once        sync.Once
}

func New(root string) *Controller {
	return &Controller{root: root}
}

func (c *Controller) Mux() *http.ServeMux {
	c.once.Do(func() {
		m := http.NewServeMux()
		m.Handle("/ws", websocket.Handler(c.WSHandler))
		c.mux = m
	})

	return c.mux
}

func (c *Controller) SetCEClient(client ceclient.Client) {
	ceClient = client
	log.Println("Getting external clusters usernames")
	manager.getExternalUsernames()
}

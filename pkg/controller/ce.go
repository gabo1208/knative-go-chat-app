package controller

import (
	"encoding/json"
	"fmt"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

const (
	ChatAppEventSource = "gabo1208.go-chat-client/source"
	NewUserConnected   = "gabo1208.go-chat-client.NewUserConnected"
	UserDisconnected   = "gabo1208.go-chat-client.UserDisconnected"
)

func (c *Controller) CeHandler(event cloudevents.Event) {
	fmt.Println("got", event.String())

	// TODO: cloudevents needs a websocket transport.

	b, err := json.Marshal(event)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	manager.broadcast <- string(b)
}

func (c *client) SendCE(ceType, contentType string, data interface{}, globalEvent bool) {
	cloudEvent := cloudevents.NewEvent()
	cloudEvent.SetID(uuid.NewString())
	cloudEvent.SetSource(ChatAppEventSource)
	cloudEvent.SetType(ceType)
	if err := cloudEvent.SetData(contentType, data); err != nil {
		log.Println(err)
	}
	websocket.JSON.Send(c.socket, cloudEvent)
}

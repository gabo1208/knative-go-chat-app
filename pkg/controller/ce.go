package controller

import (
	"fmt"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

const (
	ChatAppEventSource  = "gabo1208.knative-go-chat-app/source"
	FirstUserConnection = "gabo1208.knative-go-chat-app.FirstUserConnection"
	UserReconnected 		= "gabo1208.knative-go-chat-app.UserReconnected"
	NewUserConnected    = "gabo1208.knative-go-chat-app.NewUserConnected"
	UserDisconnected    = "gabo1208.knative-go-chat-app.UserDisconnected"
)

func (c *Controller) CeHandler(event cloudevents.Event) {
	fmt.Println("got", event.String())
	manager.broadcast <- &event
}

func (c *client) createCE(ceType, contentType string, data interface{}) *cloudevents.Event {
	cloudEvent := cloudevents.NewEvent()
	cloudEvent.SetID(uuid.NewString())
	cloudEvent.SetSource(ChatAppEventSource)
	cloudEvent.SetType(ceType)
	if err := cloudEvent.SetData(contentType, data); err != nil {
		log.Println(err)
		return nil
	}

	return &cloudEvent
}

func (c *client) SendCE(ce *cloudevents.Event, globalEvent bool) {
	if globalEvent {
		manager.broadcast <- ce
	} else {
		websocket.JSON.Send(c.socket, *ce)
	}
}

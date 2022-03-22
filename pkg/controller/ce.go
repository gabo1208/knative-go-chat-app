package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

const (
	ChatAppEventSource  = "gabo1208.knative-go-chat-app/source"
	FirstUserConnection = "gabo1208.knative-go-chat-app.FirstUserConnection"
	UserReconnected     = "gabo1208.knative-go-chat-app.UserReconnected"
	NewUserConnected    = "gabo1208.knative-go-chat-app.NewUserConnected"
	UserDisconnected    = "gabo1208.knative-go-chat-app.UserDisconnected"
)

func (c *Controller) CeHandler(event cloudevents.Event) {
	fmt.Println("got", event.String())

	var msg map[string]interface{}
	if err := json.Unmarshal(event.Data(), &msg); err != nil {
		log.Print(err)
		return
	}

	if value, ok := msg["to"]; !ok {
		manager.broadcast <- &event
	} else if user, ok := manager.usernames[value.(string)]; ok {
		log.Printf("receiving %s", msg)
		err := user.processWSMessage(msg)
		if err != nil {
			log.Print(err)
			return
		}
	}
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
		otherClusters := os.Getenv("CLUSTERS_BROKERS_URI")
		for _, uri := range strings.Split(otherClusters, ",") {
			log.Printf("sending to %s", uri)
			ctx := cloudevents.ContextWithTarget(context.Background(), uri)
			if result := ceClient.Send(ctx, *ce); cloudevents.IsUndelivered(result) {
				log.Fatalf("failed to send, %v", result)
			} else {
				log.Printf("res %s", result)
			}
		}
	} else {
		websocket.JSON.Send(c.socket, *ce)
	}
}

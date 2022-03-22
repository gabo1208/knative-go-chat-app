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
	MessageFromUser     = "gabo1208.knative-go-chat-app.MessageFromUser"
	GetUsers            = "gabo1208.knative-go-chat-app.GetUsers"
	ExternalUsers       = "gabo1208.knative-go-chat-app.ExternalUsers"
)

func (c *Controller) CeHandler(event cloudevents.Event) {
	fmt.Println("got", event.String())
	if event.Type() == GetUsers {
		log.Printf("sending usernames %v", manager.usernames)
		SendCEViaHTTP(
			createCE(GetUsers, cloudevents.ApplicationJSON, GetUsernames(manager.usernames)),
			string(event.Data()),
		)
		return
	} else if event.Type() == ExternalUsers {
		AppendUsernames(&event, &manager)
		return
	} else if event.Type() == NewUserConnected {
		username := string(event.Data())
		log.Printf("adding external username %s", username)
		manager.usernames[username] = &client{
			id:         uuid.New().String(),
			username:   username,
			registered: true,
		}
	} else if event.Type() == UserDisconnected {
		username := string(event.Data())
		log.Printf("removing external username %s", username)
		delete(manager.usernames, username)
	}

	var msg map[string]interface{}
	if err := json.Unmarshal(event.Data(), &msg); err != nil {
		log.Print(err)
		return
	}

	if value, ok := msg["to"]; !ok {
		manager.broadcast <- &event
	} else if user, ok := manager.usernames[value.(string)]; ok {
		log.Printf("receiving %s", msg)
		err := user.processWSMessage(msg, false)
		if err != nil {
			log.Print(err)
		}
	}
}

func createCE(ceType, contentType string, data interface{}) *cloudevents.Event {
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

func (c *client) SendCE(ce *cloudevents.Event, globalEvent, local bool) {
	if globalEvent {
		// If the msg should be sent to all users in the local instance
		if local {
			manager.broadcast <- ce
		}
		// Send the message to all the registered brokers
		otherClusters := os.Getenv("CLUSTERS_BROKERS_URI")
		for _, uri := range strings.Split(otherClusters, ",") {
			log.Printf("sending to %s", uri)
			SendCEViaHTTP(ce, uri)
		}
	} else {
		websocket.JSON.Send(c.socket, *ce)
	}
}

func SendCEViaHTTP(ce *cloudevents.Event, uri string) {
	ctx := cloudevents.ContextWithTarget(context.Background(), uri)
	if result := ceClient.Send(ctx, *ce); cloudevents.IsUndelivered(result) {
		log.Fatalf("failed to send, %v", result)
	} else {
		log.Printf("ce res %v", result)
	}
}

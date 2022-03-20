package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"

	"knative.dev/eventing/pkg/kncloudevents"
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
		otherClusters := os.Getenv("CLUSTERS_BROKERS_URI")
		for _, uri := range strings.Split(otherClusters, ",") {
			ctx := cloudevents.ContextWithTarget(context.TODO(), uri)
			cloudevents.ContextWithRetriesExponentialBackoff(ctx, time.Duration(500), 5)

			response, result := ceClient.Request(ctx, *ce)
			if !isSuccess(ctx, result) {
				log.Printf("Failed to deliver to %q %s", uri, response)
				return
			}
		}
	} else {
		websocket.JSON.Send(c.socket, *ce)
	}
}

func isSuccess(ctx context.Context, result protocol.Result) bool {
	var retriesResult *cehttp.RetriesResult
	if cloudevents.ResultAs(result, &retriesResult) {
		var httpResult *cehttp.Result
		if cloudevents.ResultAs(retriesResult.Result, &httpResult) {
			retry, _ := kncloudevents.SelectiveRetry(ctx, &http.Response{StatusCode: httpResult.StatusCode}, nil)
			return !retry
		}
		log.Printf("Invalid result type, not HTTP Result: %v", retriesResult.Result)
		return false
	}

	log.Printf("Invalid result type, not RetriesResult")
	return false
}

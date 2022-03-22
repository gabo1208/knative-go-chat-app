package main

import (
	"context"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func main() {
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	// Create an Event.
	event := cloudevents.NewEvent()
	event.SetSource("example/uri")
	event.SetType("example.type")
	event.SetData(cloudevents.TextPlain, "hola q tal")

	// Set a target.
	ctx := cloudevents.ContextWithTarget(context.Background(), "http://sockeye.chat-demo.127.0.0.1.sslip.io")

	// Send that Event.
	if result := c.Send(ctx, event); cloudevents.IsUndelivered(result) {
		log.Fatalf("failed to send, %v", result)
	}
}

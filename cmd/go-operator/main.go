package main

import (
	"log"
	"os"

	"github.com/gruyaume/go-operator/internal/events"
)

func main() {
	log.Println("Starting go-operator")
	eventType, err := events.GetEventType()
	if err != nil {
		log.Println("could not get event type:", err)
		os.Exit(0)
	}

	log.Println("Event type:", eventType)

	os.Exit(0)
}

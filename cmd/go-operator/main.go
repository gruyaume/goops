package main

import (
	"log"
	"os"

	"github.com/gruyaume/go-operator/internal/commands"
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
	err = commands.SetStatus(commands.StatusActive)
	if err != nil {
		log.Println("could not set status:", err)
		os.Exit(1)
	}
	log.Println("Status set to active")
	log.Println("go-operator finished")
	os.Exit(0)
}

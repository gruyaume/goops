package events

import (
	"fmt"
	"os"
)

const JUJU_HOOK_NAME_ENV = "JUJU_HOOK_NAME"

type EventType string

const (
	EventTypeStorageAttached EventType = "storage-attached"
	EventTypeInstall         EventType = "install"
	EventTypeLeaderElected   EventType = "leader-elected"
	EventTypeConfigChanged   EventType = "config-changed"
	EventTypeStart           EventType = "start"
	EventTypeStop            EventType = "stop"
	EventTypeUpdateStatus    EventType = "update-status"
)

func GetEventType() (EventType, error) {
	hookName := os.Getenv(JUJU_HOOK_NAME_ENV)
	if hookName == "" {
		return "", fmt.Errorf("environment variable %s is not set", JUJU_HOOK_NAME_ENV)
	}
	switch hookName {
	case "storage-attached":
		return EventTypeStorageAttached, nil
	case "install":
		return EventTypeInstall, nil
	case "start":
		return EventTypeStart, nil
	case "leader-elected":
		return EventTypeLeaderElected, nil
	case "config-changed":
		return EventTypeConfigChanged, nil
	case "stop":
		return EventTypeStop, nil
	case "update-status":
		return EventTypeUpdateStatus, nil
	default:
		return "", fmt.Errorf("unknown event type: %s", hookName)
	}
}

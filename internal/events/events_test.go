package events_test

import (
	"testing"

	"github.com/gruyaume/go-operator/internal/events"
)

func TestValidEvents(t *testing.T) {

	validEvents := []string{
		"storage-attached",
		"install",
		"leader-elected",
		"config-changed",
		"start",
		"stop",
		"update-status",
	}
	for _, event := range validEvents {
		t.Run(event, func(t *testing.T) {
			t.Setenv("JUJU_HOOK_NAME", event)

			eventType, err := events.GetEventType()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if eventType != events.EventType(event) {
				t.Fatalf("expected event type %s, got %s", event, eventType)
			}
		})
	}
}

func TestInValidEvents(t *testing.T) {

	validEvents := []string{
		"",
		"invalidevent",
	}
	for _, event := range validEvents {
		t.Run(event, func(t *testing.T) {
			t.Setenv("JUJU_HOOK_NAME", event)

			eventType, err := events.GetEventType()
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if eventType != "" {
				t.Fatalf("expected empty event type, got %s", eventType)
			}
		})
	}
}

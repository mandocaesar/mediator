package redis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/mandocaesar/mediator/pkg/mediator"
)

func setupTestRedis(t *testing.T) (*redis.Client, func()) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, func() {
		client.Close()
		mr.Close()
	}
}

func TestEventStore(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	store := NewEventStore(client, DefaultConfig())

	t.Run("store and retrieve events", func(t *testing.T) {
		ctx := context.Background()
		event := mediator.Event{
			Name:    "test.event",
			Payload: map[string]interface{}{"key": "value"},
		}

		// Store event
		err := store.StoreEvent(ctx, event)
		if err != nil {
			t.Fatalf("Failed to store event: %v", err)
		}

		// Retrieve events
		events, err := store.GetEvents(ctx, "test.event", 10)
		if err != nil {
			t.Fatalf("Failed to get events: %v", err)
		}

		if len(events) != 1 {
			t.Fatalf("Expected 1 event, got %d", len(events))
		}

		if events[0]["name"] != "test.event" {
			t.Errorf("Expected event name 'test.event', got %v", events[0]["name"])
		}
	})

	t.Run("clear events", func(t *testing.T) {
		ctx := context.Background()
		event := mediator.Event{
			Name:    "clear.test",
			Payload: map[string]interface{}{"key": "value"},
		}

		// Store event
		err := store.StoreEvent(ctx, event)
		if err != nil {
			t.Fatalf("Failed to store event: %v", err)
		}

		// Clear events
		err = store.ClearEvents(ctx, "clear.test")
		if err != nil {
			t.Fatalf("Failed to clear events: %v", err)
		}

		// Verify events are cleared
		events, err := store.GetEvents(ctx, "clear.test", 10)
		if err != nil {
			t.Fatalf("Failed to get events: %v", err)
		}

		if len(events) != 0 {
			t.Errorf("Expected 0 events after clear, got %d", len(events))
		}
	})
}

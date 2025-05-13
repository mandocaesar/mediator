package mediator

import "context"

// EventStore defines the interface for event storage
type EventStore interface {
	// StoreEvent stores an event
	StoreEvent(ctx context.Context, event Event) error

	// GetEvents retrieves events by event name
	GetEvents(ctx context.Context, eventName string, limit int64) ([]map[string]interface{}, error)

	// ClearEvents removes all events for a given event name
	ClearEvents(ctx context.Context, eventName string) error
}

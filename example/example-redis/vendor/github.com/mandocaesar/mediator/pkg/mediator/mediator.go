package mediator

import (
	"context"
	"fmt"
	"sync"
)

// Event represents a generic event in the system
type Event struct {
	Name    string
	Payload interface{}
}

// Mediator manages event subscriptions and publishing
type Mediator struct {
	subscribers map[string][]EventHandler
	eventStore  EventStore
	mu          sync.RWMutex
}

// EventHandler is a function type that handles events
type EventHandler func(ctx context.Context, event Event) error

var (
	globalMediator *Mediator
	mediatorOnce   sync.Once
)

// New creates a singleton Mediator instance
func New() *Mediator {
	mediatorOnce.Do(func() {
		globalMediator = &Mediator{
			subscribers: make(map[string][]EventHandler),
		}
	})
	return globalMediator
}

// SetEventStore sets the event store for the mediator
func (m *Mediator) SetEventStore(store EventStore) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.eventStore = store
}

// GetMediator returns the existing mediator instance
func GetMediator() *Mediator {
	if globalMediator == nil {
		return New()
	}
	return globalMediator
}

// Subscribe adds an event handler for a specific event type
func (m *Mediator) Subscribe(eventName string, handler EventHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscribers[eventName] = append(m.subscribers[eventName], handler)
}

// Publish sends an event to all registered handlers and stores it if event store is configured
func (m *Mediator) Publish(ctx context.Context, event Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	handlers, exists := m.subscribers[event.Name]
	if !exists {
		return fmt.Errorf("no handlers for event: %s", event.Name)
	}

	var errs []error
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			errs = append(errs, err)
		}
	}

	// Store event if event store is configured
	if m.eventStore != nil {
		if err := m.eventStore.StoreEvent(ctx, event); err != nil {
			errs = append(errs, fmt.Errorf("failed to store event: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors in event handlers: %v", errs)
	}

	return nil
}

// GetEvents retrieves events from the event store
func (m *Mediator) GetEvents(ctx context.Context, eventName string, limit int64) ([]map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.eventStore == nil {
		return nil, fmt.Errorf("no event store configured")
	}

	return m.eventStore.GetEvents(ctx, eventName, limit)
}

// ClearEvents removes all events for a given event name
func (m *Mediator) ClearEvents(ctx context.Context, eventName string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.eventStore == nil {
		return fmt.Errorf("no event store configured")
	}

	return m.eventStore.ClearEvents(ctx, eventName)
}

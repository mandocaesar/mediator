package mediator

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	// Reset global mediator and once for testing
	globalMediator = nil
	mediatorOnce = sync.Once{}

	// Test creating new instance
	m1 := New()
	if m1 == nil {
		t.Error("New() returned nil")
	}

	// Test singleton pattern
	m2 := New()
	if m1 != m2 {
		t.Error("New() did not return singleton instance")
	}
}

func TestGetMediator(t *testing.T) {
	// Reset global mediator and once for testing
	globalMediator = nil
	mediatorOnce = sync.Once{}

	// Test getting instance when none exists
	m1 := GetMediator()
	if m1 == nil {
		t.Error("GetMediator() returned nil")
	}

	// Test getting existing instance
	m2 := GetMediator()
	if m1 != m2 {
		t.Error("GetMediator() did not return singleton instance")
	}
}

func TestMediator_Subscribe(t *testing.T) {
	m := &Mediator{
		subscribers: make(map[string][]EventHandler),
	}

	eventName := "test.event"
	handler := func(ctx context.Context, event Event) error { return nil }

	// Test subscribing single handler
	m.Subscribe(eventName, handler)
	if len(m.subscribers[eventName]) != 1 {
		t.Errorf("Subscribe() failed to add handler, got %d handlers", len(m.subscribers[eventName]))
	}

	// Test subscribing multiple handlers
	m.Subscribe(eventName, handler)
	if len(m.subscribers[eventName]) != 2 {
		t.Errorf("Subscribe() failed to add multiple handlers, got %d handlers", len(m.subscribers[eventName]))
	}
}

func TestMediator_Publish(t *testing.T) {
	tests := []struct {
		name       string
		eventName  string
		setupMock  func() *Mediator
		wantErr    bool
		errMessage string
	}{
		{
			name:      "successful publish",
			eventName: "test.success",
			setupMock: func() *Mediator {
				m := &Mediator{
					subscribers: make(map[string][]EventHandler),
				}
				m.Subscribe("test.success", func(ctx context.Context, event Event) error {
					return nil
				})
				return m
			},
			wantErr: false,
		},
		{
			name:      "no handlers",
			eventName: "test.nohandlers",
			setupMock: func() *Mediator {
				return &Mediator{
					subscribers: make(map[string][]EventHandler),
				}
			},
			wantErr:    true,
			errMessage: "no handlers for event: test.nohandlers",
		},
		{
			name:      "handler error",
			eventName: "test.error",
			setupMock: func() *Mediator {
				m := &Mediator{
					subscribers: make(map[string][]EventHandler),
				}
				m.Subscribe("test.error", func(ctx context.Context, event Event) error {
					return errors.New("handler error")
				})
				return m
			},
			wantErr:    true,
			errMessage: "errors in event handlers",
		},
		{
			name:      "multiple handlers with error",
			eventName: "test.multiple",
			setupMock: func() *Mediator {
				m := &Mediator{
					subscribers: make(map[string][]EventHandler),
				}
				m.Subscribe("test.multiple", func(ctx context.Context, event Event) error {
					return nil
				})
				m.Subscribe("test.multiple", func(ctx context.Context, event Event) error {
					return errors.New("second handler error")
				})
				return m
			},
			wantErr:    true,
			errMessage: "errors in event handlers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupMock()
			ctx := context.Background()
			event := Event{
				Name:    tt.eventName,
				Payload: "test payload",
			}

			err := m.Publish(ctx, event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMessage) {
				t.Errorf("Publish() error message = %v, want %v", err, tt.errMessage)
			}
		})
	}
}

package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mandocaesar/mediator/pkg/mediator"
)

func TestEventStore(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Set up expectations for table creation
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE INDEX IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE INDEX IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(0, 0))

	// Create a new event store
	store, err := NewEventStore(db, DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create event store: %v", err)
	}

	t.Run("store and retrieve events", func(t *testing.T) {
		ctx := context.Background()
		event := mediator.Event{
			Name:    "test.event",
			Payload: map[string]interface{}{"key": "value"},
		}

		// Expect the insert query to be executed
		mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))

		// Expect the trim query to be executed
		mock.ExpectExec("DELETE FROM").WillReturnResult(sqlmock.NewResult(0, 0))

		// Store the event
		err := store.StoreEvent(ctx, event)
		if err != nil {
			t.Fatalf("Failed to store event: %v", err)
		}

		// Expect the select query to be executed
		rows := sqlmock.NewRows([]string{"event_data"}).
			AddRow(`{"name":"test.event","payload":{"key":"value"},"timestamp":"2025-05-11T13:00:00Z"}`)
		mock.ExpectQuery("SELECT event_data").WillReturnRows(rows)

		// Get events
		events, err := store.GetEvents(ctx, "test.event", 10)
		if err != nil {
			t.Fatalf("Failed to get events: %v", err)
		}

		// Check that we got the expected event
		if len(events) != 1 {
			t.Fatalf("Expected 1 event, got %d", len(events))
		}

		payload, ok := events[0]["payload"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to get payload")
		}

		if payload["key"] != "value" {
			t.Errorf("Expected payload key 'value', got '%v'", payload["key"])
		}
	})

	t.Run("clear events", func(t *testing.T) {
		ctx := context.Background()

		// Expect the delete query to be executed
		mock.ExpectExec("DELETE FROM").WillReturnResult(sqlmock.NewResult(0, 2))

		// Clear events
		err := store.ClearEvents(ctx, "test.event")
		if err != nil {
			t.Fatalf("Failed to clear events: %v", err)
		}

		// Expect the select query to return no rows
		rows := sqlmock.NewRows([]string{"event_data"})
		mock.ExpectQuery("SELECT event_data").WillReturnRows(rows)

		// Get events after clearing
		events, err := store.GetEvents(ctx, "test.event", 10)
		if err != nil {
			t.Fatalf("Failed to get events: %v", err)
		}

		// Check that we got no events
		if len(events) != 0 {
			t.Errorf("Expected 0 events, got %d", len(events))
		}
	})

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

// TestWithRealDB is a more comprehensive test using a real database connection
// This test is skipped by default and can be enabled by setting the POSTGRES_TEST_DSN environment variable
func TestWithRealDB(t *testing.T) {
	// Skip this test unless explicitly enabled
	if os.Getenv("POSTGRES_TEST_DSN") == "" {
		t.Skip("Skipping PostgreSQL integration test. Set POSTGRES_TEST_DSN to enable.")
	}
	dsn := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	// Connect to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a unique table name for this test run
	config := DefaultConfig()
	config.Prefix = "mediator_events_test"

	// Create a new event store
	store, err := NewEventStore(db, config)
	if err != nil {
		t.Fatalf("Failed to create event store: %v", err)
	}

	// Clean up before and after the test
	store.ClearEvents(context.Background(), "test.event")
	defer store.ClearEvents(context.Background(), "test.event")

	// Test storing and retrieving events
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		event := mediator.Event{
			Name:    "test.event",
			Payload: map[string]interface{}{"index": i},
		}
		if err := store.StoreEvent(ctx, event); err != nil {
			t.Fatalf("Failed to store event: %v", err)
		}
	}

	// Get events with limit
	events, err := store.GetEvents(ctx, "test.event", 3)
	if err != nil {
		t.Fatalf("Failed to get events: %v", err)
	}

	// Check that we got the expected number of events
	if len(events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(events))
	}

	// Check that we got the most recent events (highest indices)
	for i, event := range events {
		payload := event["payload"].(map[string]interface{})
		index := int(payload["index"].(float64))
		expectedIndex := 4 - i // We expect indices 4, 3, 2 in that order
		if index != expectedIndex {
			t.Errorf("Event at position %d: expected index %d, got %d", i, expectedIndex, index)
		}
	}

	// Test clearing events
	if err := store.ClearEvents(ctx, "test.event"); err != nil {
		t.Fatalf("Failed to clear events: %v", err)
	}

	// Check that all events were cleared
	events, err = store.GetEvents(ctx, "test.event", 10)
	if err != nil {
		t.Fatalf("Failed to get events: %v", err)
	}

	if len(events) != 0 {
		t.Errorf("Expected 0 events after clearing, got %d", len(events))
	}
}

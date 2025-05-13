package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/mandocaesar/mediator/pkg/mediator"
)

// EventStore represents a PostgreSQL-based event store
type EventStore struct {
	db     *sql.DB
	prefix string
}

// Config represents PostgreSQL event store configuration
type Config struct {
	Prefix           string
	MaxEventsPerType int64
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		Prefix:           "mediator_events",
		MaxEventsPerType: 1000,
	}
}

// NewEventStore creates a new PostgreSQL event store
func NewEventStore(db *sql.DB, config Config) (*EventStore, error) {
	if config.Prefix == "" {
		config.Prefix = DefaultConfig().Prefix
	}

	store := &EventStore{
		db:     db,
		prefix: config.Prefix,
	}

	// Initialize tables
	if err := store.initTables(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return store, nil
}

// initTables creates the necessary tables if they don't exist
func (s *EventStore) initTables(ctx context.Context) error {
	// Create events table
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			event_name TEXT NOT NULL,
			event_data JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)
	`, pq.QuoteIdentifier(s.prefix))

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create events table: %w", err)
	}

	// Create index on event_name for faster lookups
	indexQuery := fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS %s_event_name_idx ON %s (event_name)
	`, s.prefix, pq.QuoteIdentifier(s.prefix))

	_, err = s.db.ExecContext(ctx, indexQuery)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// Create index on created_at for faster sorting
	timeIndexQuery := fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS %s_created_at_idx ON %s (created_at)
	`, s.prefix, pq.QuoteIdentifier(s.prefix))

	_, err = s.db.ExecContext(ctx, timeIndexQuery)
	if err != nil {
		return fmt.Errorf("failed to create time index: %w", err)
	}

	return nil
}

// StoreEvent stores an event in PostgreSQL
func (s *EventStore) StoreEvent(ctx context.Context, event mediator.Event) error {
	// Create event data with metadata
	timestamp := time.Now().UTC()
	eventData := map[string]interface{}{
		"name":      event.Name,
		"payload":   event.Payload,
		"timestamp": timestamp,
	}

	// Convert to JSON
	data, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Insert event
	query := fmt.Sprintf(`
		INSERT INTO %s (event_name, event_data, created_at)
		VALUES ($1, $2, $3)
	`, pq.QuoteIdentifier(s.prefix))

	_, err = s.db.ExecContext(ctx, query, event.Name, data, timestamp)
	if err != nil {
		return fmt.Errorf("failed to store event: %w", err)
	}

	// Trim events if needed
	if DefaultConfig().MaxEventsPerType > 0 {
		err = s.trimEvents(ctx, event.Name)
		if err != nil {
			return fmt.Errorf("failed to trim events: %w", err)
		}
	}

	return nil
}

// trimEvents ensures that only the most recent MaxEventsPerType events are kept
func (s *EventStore) trimEvents(ctx context.Context, eventName string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id IN (
			SELECT id FROM %s
			WHERE event_name = $1
			ORDER BY created_at DESC
			OFFSET $2
		)
	`, pq.QuoteIdentifier(s.prefix), pq.QuoteIdentifier(s.prefix))

	_, err := s.db.ExecContext(ctx, query, eventName, DefaultConfig().MaxEventsPerType)
	if err != nil {
		return fmt.Errorf("failed to trim events: %w", err)
	}

	return nil
}

// GetEvents retrieves events from PostgreSQL by event name
func (s *EventStore) GetEvents(ctx context.Context, eventName string, limit int64) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = DefaultConfig().MaxEventsPerType
	}

	// Query for events
	query := fmt.Sprintf(`
		SELECT event_data
		FROM %s
		WHERE event_name = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, pq.QuoteIdentifier(s.prefix))

	rows, err := s.db.QueryContext(ctx, query, eventName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	events := make([]map[string]interface{}, 0)
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan event data: %w", err)
		}

		var event map[string]interface{}
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}

// ClearEvents removes all events for a given event name
func (s *EventStore) ClearEvents(ctx context.Context, eventName string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE event_name = $1
	`, pq.QuoteIdentifier(s.prefix))

	_, err := s.db.ExecContext(ctx, query, eventName)
	if err != nil {
		return fmt.Errorf("failed to clear events: %w", err)
	}

	return nil
}

// Close closes the database connection
func (s *EventStore) Close() error {
	return s.db.Close()
}

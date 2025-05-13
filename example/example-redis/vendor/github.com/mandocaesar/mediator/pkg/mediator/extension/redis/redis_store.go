package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mandocaesar/mediator/pkg/mediator"
)

// EventStore represents a Redis-based event store
type EventStore struct {
	client *redis.Client
	prefix string
}

// Config represents Redis event store configuration
type Config struct {
	Prefix           string
	EventTTL         time.Duration
	MaxEventsPerType int64
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		Prefix:           "mediator:events",
		EventTTL:         24 * time.Hour,
		MaxEventsPerType: 1000,
	}
}

// NewEventStore creates a new Redis event store
func NewEventStore(client *redis.Client, config Config) *EventStore {
	if config.Prefix == "" {
		config.Prefix = DefaultConfig().Prefix
	}
	return &EventStore{
		client: client,
		prefix: config.Prefix,
	}
}

// StoreEvent stores an event in Redis
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

	// Generate key with timestamp for ordering
	key := fmt.Sprintf("%s:%s:%d", s.prefix, event.Name, timestamp.UnixNano())

	// Store event with expiration
	err = s.client.Set(ctx, key, data, DefaultConfig().EventTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to store event: %w", err)
	}

	// Add to time series list
	listKey := fmt.Sprintf("%s:%s:timeline", s.prefix, event.Name)
	err = s.client.RPush(ctx, listKey, key).Err()
	if err != nil {
		return fmt.Errorf("failed to push event to list: %w", err)
	}

	return nil
}

// GetEvents retrieves events from Redis by event name
func (s *EventStore) GetEvents(ctx context.Context, eventName string, limit int64) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = DefaultConfig().MaxEventsPerType
	}

	// Get event keys from timeline
	listKey := fmt.Sprintf("%s:%s:timeline", s.prefix, eventName)
	// Get most recent events
	keys, err := s.client.LRange(ctx, listKey, -limit, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get event keys: %w", err)
	}

	if len(keys) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Get events data
	pipe := s.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	// Process results
	events := make([]map[string]interface{}, 0, len(cmds))
	for _, cmd := range cmds {
		data, err := cmd.Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get event data: %w", err)
		}

		var event map[string]interface{}
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

// ClearEvents removes all events for a given event name
func (s *EventStore) ClearEvents(ctx context.Context, eventName string) error {
	// Get event keys from timeline
	listKey := fmt.Sprintf("%s:%s:timeline", s.prefix, eventName)
	keys, err := s.client.LRange(ctx, listKey, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get event keys: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	// Delete all events and timeline
	pipe := s.client.Pipeline()
	for _, key := range keys {
		pipe.Del(ctx, key)
	}
	pipe.Del(ctx, listKey)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to clear events: %w", err)
	}

	return nil
}

// Close closes the Redis client
func (s *EventStore) Close() error {
	return s.client.Close()
}

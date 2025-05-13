# Redis Event Store for Mediator

This extension provides a Redis implementation of the `EventStore` interface for the mediator library, allowing you to store and retrieve events using Redis.

## Features

- Store events in Redis with automatic expiration
- Retrieve events by name with optional limits
- Clear events by name
- Chronological event ordering
- Configurable event TTL

## Installation

```bash
go get github.com/yourusername/mediator
```

## Dependencies

This extension requires the following dependencies:

```bash
go get github.com/go-redis/redis/v8
```

## Usage

### Basic Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/yourusername/mediator/pkg/mediator"
	redisstore "github.com/yourusername/mediator/pkg/mediator/extension/redis"
)

func main() {
	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

	// Create event store
	config := redisstore.DefaultConfig()
	config.Prefix = "my_events" // Key prefix
	config.EventTTL = 24 * time.Hour // Event expiration time

	store := redisstore.NewEventStore(client, config)
	defer store.Close()

	// Create mediator
	m := mediator.New()

	// Register event store
	m.RegisterEventStore(store)

	// Use mediator to publish events
	ctx := context.Background()
	event := mediator.Event{
		Name:    "user.created",
		Payload: map[string]interface{}{"id": 123, "name": "John Doe"},
	}

	if err := m.Publish(ctx, event); err != nil {
		log.Fatalf("Failed to publish event: %v", err)
	}

	// Retrieve events
	events, err := store.GetEvents(ctx, "user.created", 10)
	if err != nil {
		log.Fatalf("Failed to get events: %v", err)
	}

	for _, event := range events {
		fmt.Printf("Event: %+v\n", event)
	}
}
```

### Configuration Options

The Redis event store can be configured with the following options:

- `Prefix`: The key prefix for Redis keys (default: "mediator:events")
- `EventTTL`: Time-to-live for events (default: 24 hours)
- `MaxEventsPerType`: Maximum number of events to keep per event type (default: 1000)

## Redis Data Structure

The extension uses the following Redis data structures:

- **Keys**: `{prefix}:{event_name}:{timestamp}` - Stores the event data as JSON
- **Lists**: `{prefix}:{event_name}:timeline` - Stores the keys of events in chronological order

## Event Retrieval

Events are retrieved in reverse chronological order (newest first) using Redis' `LRANGE` command with negative indices. This ensures that you always get the most recent events when using limits.

## Testing

The extension includes tests using a mock Redis server (miniredis). To run the tests:

```bash
go test -v ./pkg/mediator/extension/redis/...
```

## Example with Custom Redis Options

```go
// Connect to Redis with custom options
client := redis.NewClient(&redis.Options{
    Addr:     "redis.example.com:6379",
    Password: "password123",
    DB:       1,
    PoolSize: 10,
})

// Create event store with custom configuration
config := redisstore.DefaultConfig()
config.Prefix = "app_events"
config.EventTTL = 7 * 24 * time.Hour // 1 week

store := redisstore.NewEventStore(client, config)
```

## License

This project is licensed under the same license as the mediator library.

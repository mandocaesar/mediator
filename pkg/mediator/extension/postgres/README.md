# PostgreSQL Event Store for Mediator

This extension provides a PostgreSQL implementation of the `EventStore` interface for the mediator library, allowing you to store and retrieve events using a PostgreSQL database.

## Features

- Store events in a PostgreSQL database
- Retrieve events by name with optional limits
- Clear events by name
- Automatic table and index creation
- Configurable event limit per event type

## Installation

```bash
go get github.com/yourusername/mediator
```

## Dependencies

This extension requires the following dependencies:

```bash
go get github.com/lib/pq
```

## Usage

### Basic Usage

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/yourusername/mediator/pkg/mediator"
	"github.com/yourusername/mediator/pkg/mediator/extension/postgres"
)

func main() {
	// Connect to PostgreSQL
	db, err := sql.Open("postgres", "postgres://username:password@localhost:5432/database?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create event store
	config := postgres.DefaultConfig()
	config.Prefix = "my_events" // Table name prefix
	config.MaxEventsPerType = 100 // Maximum events to keep per event type

	store, err := postgres.NewEventStore(db, config)
	if err != nil {
		log.Fatalf("Failed to create event store: %v", err)
	}

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

The PostgreSQL event store can be configured with the following options:

- `Prefix`: The table name prefix (default: "mediator_events")
- `MaxEventsPerType`: Maximum number of events to keep per event type (default: 1000)

## Database Schema

The extension creates the following database objects:

- A table named `{prefix}` with columns:
  - `id`: Serial primary key
  - `event_name`: Text, the name of the event
  - `event_data`: JSONB, the event data including payload and metadata
  - `created_at`: Timestamp with timezone, when the event was created

- Indexes:
  - `{prefix}_event_name_idx`: Index on `event_name` for faster lookups
  - `{prefix}_created_at_idx`: Index on `created_at` for faster sorting

## Event Trimming

The PostgreSQL event store automatically trims events when the number of events for a specific event type exceeds the configured `MaxEventsPerType`. Only the most recent events are kept, based on their creation timestamp.

## Testing

The extension includes both unit tests using a mock database and integration tests using a real PostgreSQL database. To run the integration tests, you need to have a PostgreSQL database available and set the connection string in the test file.

```bash
go test -v ./pkg/mediator/extension/postgres/...
```

## License

This project is licensed under the same license as the mediator library.

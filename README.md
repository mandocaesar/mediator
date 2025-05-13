# Go Mediator Pattern Implementation

A lightweight, thread-safe mediator library for Go that implements the mediator pattern with event-driven communication. This library provides a robust solution for decoupling components in your application through event-based interactions.

## Features

- 🔒 Thread-safe event publishing and subscription
- 🌟 Singleton mediator pattern for global access
- ⚡ Asynchronous event handling
- 🔄 Multiple event store implementations (Redis, PostgreSQL)
- 📦 Easy-to-use API
- 🧪 High test coverage
- 📝 Comprehensive documentation

## Installation

```bash
go get github.com/mandocaesar/mediator
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/mandocaesar/mediator/pkg/mediator"
)

func main() {
    // Get mediator instance
    m := mediator.GetMediator()

    // Subscribe to events
    m.Subscribe("user.created", func(ctx context.Context, event mediator.Event) error {
        // Handle the event
        return nil
    })

    // Publish an event
    event := mediator.Event{
        Name:    "user.created",
        Payload: map[string]interface{}{"id": "123", "name": "John Doe"},
    }
    m.Publish(context.Background(), event)
}
```

## Event Store Support

### Redis Event Store

```go
import (
    "github.com/go-redis/redis/v8"
    "github.com/mandocaesar/mediator/pkg/mediator"
    redisstore "github.com/mandocaesar/mediator/pkg/mediator/extension/redis"
)

// Create Redis client
client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Create event store
store := redisstore.NewEventStore(client, redisstore.DefaultConfig())

// Set event store in mediator
m := mediator.GetMediator()
m.SetEventStore(store)
```

### PostgreSQL Event Store

```go
import (
    "database/sql"
    "github.com/mandocaesar/mediator/pkg/mediator"
    postgresstore "github.com/mandocaesar/mediator/pkg/mediator/extension/postgres"
)

// Create PostgreSQL connection
db, _ := sql.Open("postgres", "postgres://user:pass@localhost:5432/dbname")

// Create event store
store, _ := postgresstore.NewEventStore(db, postgresstore.DefaultConfig())

// Set event store in mediator
m := mediator.GetMediator()
m.SetEventStore(store)
```

## Project Structure

```
├── pkg/
│   └── mediator/           # Core mediator package
│       ├── mediator.go     # Main mediator implementation
│       ├── event_store.go  # Event storage interface
│       └── extension/      # Event store implementations
│           ├── redis/      # Redis event store
│           └── postgres/   # PostgreSQL event store
└── example/               # Example implementations
    ├── example-app/      # Full application example
    ├── example-redis/    # Redis example
    └── example-postgres/ # PostgreSQL example
```

## Examples

The repository includes several examples:

1. **Full Application Example** (`example/example-app/`): Demonstrates a complete application using the mediator pattern with domain-driven design.
2. **Redis Example** (`example/example-redis/`): Shows how to use the mediator with Redis event store.
3. **PostgreSQL Example** (`example/example-postgres/`): Shows how to use the mediator with PostgreSQL event store.

## Running Examples

Each example includes a Docker Compose file for easy setup:

```bash
# Redis Example
cd example/example-redis
docker-compose up

# PostgreSQL Example
cd example/example-postgres
docker-compose up

# Full Application Example
cd example/example-app
docker-compose up
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

- `product.created`: Triggered when a new product is created
- `product.updated`: Triggered when a product is updated
- `sku.created`: Triggered when a new SKU is created
- `sku.updated`: Triggered when an SKU is updated

## Error Handling
The library provides comprehensive error handling:

```go
// Publishing with error handling
err := mediator.Publish(ctx, event)
if err != nil {
    // Handle error
}

// Subscription with error handling
med.Subscribe("event.name", func(ctx context.Context, event mediator.Event) error {
    if err := processEvent(event); err != nil {
        return fmt.Errorf("failed to process event: %w", err)
    }
    return nil
})
```

## Redis Extension
The library includes a Redis extension for event persistence:

```go
// Initialize Redis client
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Create Redis event store with default config
config := redis.DefaultConfig()
config.Prefix = "myapp:events"
config.EventTTL = 48 * time.Hour
config.MaxEventsPerType = 1000

eventStore := redis.NewEventStore(redisClient, config)

// Configure mediator to use Redis event store
med := mediator.GetMediator()
med.SetEventStore(eventStore)

// Events will now be automatically stored in Redis
// Retrieve stored events
events, err := med.GetEvents(ctx, "event.name", 10)
if err != nil {
    // Handle error
}

// Clear stored events
err = med.ClearEvents(ctx, "event.name")
if err != nil {
    // Handle error
}
```

The Redis extension provides:
- Automatic event persistence
- Configurable event TTL
- Maximum events per event type
- Event retrieval by type
- Event cleanup

## Contributing
1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License
MIT License

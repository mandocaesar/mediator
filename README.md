# Golang Mediator Library

## Overview
A lightweight, thread-safe mediator library for Golang that implements an event-driven communication pattern. This library enables loose coupling between components through event-based interactions, making it ideal for building scalable and maintainable applications.

## Features
- Thread-safe event publishing and subscription
- Singleton mediator pattern for global access
- Context support for cancellation and timeouts
- Flexible event payload system
- Comprehensive error handling
- Event persistence with Redis extension
- 100% test coverage

## Project Structure
```
├── pkg/
│   └── mediator/         # Core mediator package
│       ├── mediator.go   # Main mediator implementation
│       ├── mediator_test.go
│       ├── event_store.go # Event storage interface
│       └── extension/    # Extensions
│           └── redis/    # Redis event storage
├── example/             # Example implementation
│   ├── domain/         # Domain models
│   │   ├── product/    # Product domain
│   │   └── sku/       # SKU domain
│   ├── repository/    # Data access layer
│   ├── usecase/      # Business logic
│   └── main.go       # Example application
└── Makefile          # Build and test automation
```

## Installation
```bash
go get github.com/yourusername/mediator
```

## Usage

### Basic Usage
```go
// Get the global mediator instance
med := mediator.GetMediator()

// Subscribe to events
med.Subscribe("product.created", func(ctx context.Context, event mediator.Event) error {
    // Handle event
    return nil
})

// Publish events
med.Publish(ctx, mediator.Event{
    Name:    "product.created",
    Payload: product,
})
```

### Example Implementation
```go
// Create use cases with repositories
productUseCase := usecase.NewProductUseCase(
    productRepo,
    productDetailRepo,
)

// Handle product creation
product, err := productUseCase.CreateProduct(
    ctx,
    "Sample Product",
    "Product description",
    99.99,
)
```

## Testing
The project includes comprehensive tests with 100% coverage. Use the following make targets:

```bash
# Run all tests
make test

# Run tests with verbose output
make test-v

# Test specific package
make test-pkg pkg=./pkg/mediator

# Run tests with coverage report
make test-cover

# Run usecase-specific coverage
make test-cover-usecase

# Clean test cache
make test-clean
```

## Event Types
The library supports various event types:

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

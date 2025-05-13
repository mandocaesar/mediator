package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mandocaesar/mediator/pkg/mediator"
	redisstore "github.com/mandocaesar/mediator/pkg/mediator/extension/redis"
)

func main() {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	defer client.Close()

	store := redisstore.NewEventStore(client, redisstore.DefaultConfig())
	m := mediator.New()
	m.SetEventStore(store)

	event := mediator.Event{
		Name:    "demo.event",
		Payload: map[string]interface{}{"msg": "Hello from Redis!", "ts": time.Now().Format(time.RFC3339)},
	}
	if err := m.Publish(ctx, event); err != nil {
		log.Fatalf("Failed to publish event: %v", err)
	}

	fmt.Println("Events from Redis:")
	events, _ := store.GetEvents(ctx, "demo.event", 10)
	for _, e := range events {
		fmt.Printf("  %+v\n", e)
	}
}

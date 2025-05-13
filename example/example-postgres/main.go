package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/mandocaesar/mediator/pkg/mediator"
	postgresstore "github.com/mandocaesar/mediator/pkg/mediator/extension/postgres"
)

func main() {
	ctx := context.Background()
	pgDSN := "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
	pgDB, err := sql.Open("postgres", pgDSN)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer pgDB.Close()

	store, err := postgresstore.NewEventStore(pgDB, postgresstore.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to create Postgres event store: %v", err)
	}

	m := mediator.New()
	m.SetEventStore(store)

	event := mediator.Event{
		Name:    "demo.event",
		Payload: map[string]interface{}{"msg": "Hello from Postgres!", "ts": time.Now().Format(time.RFC3339)},
	}
	if err := m.Publish(ctx, event); err != nil {
		log.Fatalf("Failed to publish event: %v", err)
	}

	fmt.Println("Events from Postgres:")
	events, _ := store.GetEvents(ctx, "demo.event", 10)
	for _, e := range events {
		fmt.Printf("  %+v\n", e)
	}
}

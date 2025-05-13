module example-app

go 1.21

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/lib/pq v1.10.9
	github.com/yourusername/mediator v0.0.0
)

replace github.com/yourusername/mediator => ../../pkg/mediator

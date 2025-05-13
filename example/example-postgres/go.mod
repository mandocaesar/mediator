module example-postgres

go 1.21

require (
	github.com/lib/pq v1.10.9
	github.com/mandocaesar/mediator v0.0.0
)

replace github.com/mandocaesar/mediator => ../../pkg/mediator

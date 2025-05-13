.PHONY: test test-v test-pkg test-example test-cover test-cover-usecase test-clean

# Default test target
test:
	go test ./...

# Verbose test with test case names
test-v:
	go test -v ./...

# Test specific package with verbose output
test-pkg:
	@if [ "$(pkg)" = "" ]; then \
		echo "Usage: make test-pkg pkg=<package_path>"; \
		echo "Example: make test-pkg pkg=./pkg/mediator"; \
		exit 1; \
	fi
	go test -v $(pkg)

# Test example package with verbose output
test-example:
	go test -v ./example/...

# Run tests with coverage report for all packages
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run tests with coverage report for usecase package only
test-cover-usecase:
	go test -coverprofile=coverage.out ./example/usecase/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated for usecase package"
	@go tool cover -func=coverage.out

# Clean test cache and coverage files
test-clean:
	go clean -testcache
	rm -f coverage.out coverage.html

# Help target
help:
	@echo "Available targets:"
	@echo "  test             - Run all tests"
	@echo "  test-v           - Run all tests with verbose output"
	@echo "  test-pkg         - Test specific package (usage: make test-pkg pkg=./pkg/mediator)"
	@echo "  test-example     - Test example package"
	@echo "  test-cover       - Run tests with coverage report for all packages"
	@echo "  test-cover-usecase - Run tests with coverage report for usecase package only"
	@echo "  test-clean       - Clean test cache and coverage files"
	@echo "  help             - Show this help message"

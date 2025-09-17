.PHONY: build run-demo run-postgres clean test install deps

# Build the applications
build:
	go build -o fiber-app main.go
	go build -o fiber-demo demo.go

# Run the demo version (in-memory storage)
run-demo:
	go run demo.go

# Run the PostgreSQL version
run-postgres:
	go run main.go

# Install dependencies
deps:
	go mod tidy
	go mod download

# Clean build artifacts
clean:
	rm -f fiber-app fiber-demo

# Test the application
test:
	go test -v ./...

# Install the application
install: build
	sudo cp fiber-app /usr/local/bin/
	sudo cp fiber-demo /usr/local/bin/

# Development: run with hot reload (requires air)
dev:
	air

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Show help
help:
	@echo "Available commands:"
	@echo "  build         - Build both applications"
	@echo "  run-demo      - Run demo version with in-memory storage"
	@echo "  run-postgres  - Run PostgreSQL version"
	@echo "  deps          - Install dependencies"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  install       - Install binaries to /usr/local/bin"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  help          - Show this help"
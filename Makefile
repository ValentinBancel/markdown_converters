# Makefile for Markdown Converters

.PHONY: help build up down logs clean validate demo test

# Default target
help:
	@echo "Markdown Converters - Available commands:"
	@echo "  build    - Build all Docker services"
	@echo "  up       - Start all services"
	@echo "  down     - Stop all services"
	@echo "  logs     - View logs from all services"
	@echo "  clean    - Remove all containers and volumes"
	@echo "  validate - Run validation checks"
	@echo "  demo     - Run API demo (requires services to be running)"
	@echo "  test     - Build and test Go application"
	@echo "  dev      - Start services and run demo"

# Build all services
build:
	docker compose build

# Start all services
up:
	docker compose up -d

# Stop all services
down:
	docker compose down

# View logs
logs:
	docker compose logs -f

# Clean up everything
clean:
	docker compose down -v
	docker system prune -f

# Run validation
validate:
	./validate.sh

# Run demo (API must be running)
demo:
	./demo.sh

# Test Go application
test:
	go mod tidy
	go build -o /tmp/markdown-api .
	rm -f /tmp/markdown-api
	@echo "✅ Go application builds successfully"

# Development workflow
dev: build up
	@echo "Waiting for services to start..."
	@sleep 10
	@echo "Running demo..."
	./demo.sh

# Quick start for development
start: validate build up
	@echo "Services started successfully!"
	@echo "API: http://localhost:8080"
	@echo "pgAdmin: http://localhost:5050"
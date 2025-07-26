.PHONY: all build clean run test lint wire migrate seed help

# Default target
all: wire build

# Build the application
build:
	@echo "Building application..."
	go build -o bin/gic-crm cmd/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf tmp/

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Generate wire_gen.go files
wire:
	@echo "Generating wire_gen.go files..."
	wire ./pkg/injectors

# Run database migrations
migrate:
	@echo "Running migrations..."
	go run cmd/main.go migrate

# Seed database with initial data
seed:
	@echo "Seeding database..."
	go run cmd/main.go seed

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Default target, run wire and build"
	@echo "  build        - Build the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linter"
	@echo "  wire         - Generate wire_gen.go files for dependency injection"
	@echo "  migrate      - Run database migrations"
	@echo "  seed         - Seed database with initial data"
	@echo "  help         - Show this help message"
# Makefile for realtime-chat-backend

.PHONY: run build clean test help dev db-setup

# Variables
BINARY_NAME=chat-server
MAIN_PATH=cmd/server/main.go
BUILD_DIR=bin
MYSQL_USER?=root
MYSQL_DB=Messages

# Default target
help:
	@echo "Available commands:"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run with hot reload (requires air)"
	@echo "  make build          - Build the binary"
	@echo "  make build-prod     - Build optimized production binary"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make clean          - Remove binary files"
	@echo "  make deps           - Install/update dependencies"
	@echo "  make db-setup       - Initialize database schema"
	@echo "  make lint           - Run linter (requires golangci-lint)"
	@echo "  make fmt            - Format code"

# Run the application
run:
	go run $(MAIN_PATH)

# Run with hot reload (install air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Build the binary
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

# Build optimized production binary
build-prod:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Production binary built: $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	@echo "✅ Coverage report generated: coverage.out"

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out
	@echo "✅ Cleaned build artifacts"

# Install dependencies
deps:
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed"

# Setup database
db-setup:
	@echo "Setting up database..."
	mysql -u $(MYSQL_USER) -p $(MYSQL_DB) < create-tables.sql
	@echo "✅ Database schema created"

# Run the built binary
start: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Format code
fmt:
	go fmt ./...
	@echo "✅ Code formatted"

# Run linter (install: brew install golangci-lint or go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	golangci-lint run ./...

# Run all checks before commit
pre-commit: fmt lint test
	@echo "✅ All checks passed"

# Docker commands
docker-build:
	docker build -t chat-backend:latest .

docker-run:
	docker run -p 8080:8080 --env-file .env chat-backend:latest

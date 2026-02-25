.PHONY: help build run test lint fmt clean docker-build docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building deployes API..."
	@go build -o bin/deployes-api ./cmd/api

run: ## Run the application
	@echo "Running deployes API..."
	@go run ./cmd/api

test: ## Run tests
	@echo "Running tests..."
	@go test ./... -v -coverprofile=coverage.out

test-coverage: test ## Run tests with coverage report
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted!"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Cleaned!"

deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker compose build

docker-up: ## Start Docker containers
	@echo "Starting containers..."
	@docker compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping containers..."
	@docker compose down

docker-logs: ## Show Docker logs
	@docker compose logs -f

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@go run ./cmd/api & sleep 2 && kill $$!

db-reset: ## Reset database
	@echo "Resetting database..."
	@docker compose down -v
	@docker compose up -d

install-linter: ## Install golangci-lint
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0

check: lint test ## Run all checks (lint + test)

all: clean deps fmt check build ## Run all steps (clean, deps, fmt, check, build)

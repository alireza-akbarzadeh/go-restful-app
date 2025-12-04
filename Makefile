# Makefile for Go RESTful API Project
# This file contains common commands to manage your project

# Variables
BINARY_NAME=api-server
MIGRATE_BINARY=migrate-tool
MAIN_PATH=./cmd/server
MIGRATE_PATH=./cmd/migrate
MIGRATIONS_PATH=./cmd/migrate/migrations
DATABASE_PATH=./data.db

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Colors for terminal output
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: all build clean test coverage deps run migrate-up migrate-down migrate-create dev docker-build docker-run help

# Default target
all: deps build

## help: Display this help message
help:
	@echo "$(GREEN)Available commands:$(NC)"
	@echo ""
	@echo "$(YELLOW)Setup & Dependencies:$(NC)"
	@echo "  make deps          - Download and tidy Go dependencies"
	@echo "  make install-tools - Install development tools (air, swag, migrate)"
	@echo ""
	@echo "$(YELLOW)Development:$(NC)"
	@echo "  make run           - Run the application"
	@echo "  make dev           - Run with hot reload (requires air)"
	@echo "  make build         - Build the application binary"
	@echo "  make clean         - Remove build artifacts and temporary files"
	@echo ""
	@echo "$(YELLOW)Database:$(NC)"
	@echo "  make migrate-up    - Run database migrations up"
	@echo "  make migrate-down  - Run database migrations down"
	@echo "  make migrate-create NAME=<migration_name> - Create new migration files"
	@echo "  make db-reset      - Reset database (down and up migrations)"
	@echo "  make db-clean      - Remove database file"
	@echo ""
	@echo "$(YELLOW)Code Quality:$(NC)"
	@echo "  make test          - Run tests"
	@echo "  make coverage      - Run tests with coverage report"
	@echo "  make fmt           - Format Go code"
	@echo "  make vet           - Run go vet"
	@echo "  make lint          - Run golangci-lint (requires golangci-lint)"
	@echo "  make swagger       - Generate Swagger documentation"
	@echo ""
	@echo "$(YELLOW)Docker:$(NC)"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run application in Docker"
	@echo "  make docker-stop   - Stop Docker containers"
	@echo ""

## deps: Download and verify dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) verify
	$(GOMOD) tidy
	@echo "$(GREEN)✓ Dependencies ready!$(NC)"

## install-tools: Install development tools
install-tools:
	@echo "$(GREEN)Installing development tools...$(NC)"
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	@which migrate > /dev/null || (echo "Installing migrate..." && go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@echo "$(GREEN)✓ All tools installed!$(NC)"

## build: Build the application binary
build:
	@echo "$(GREEN)Building application...$(NC)"
	$(GOBUILD) -o bin/$(BINARY_NAME) -v $(MAIN_PATH)
	@echo "$(GREEN)✓ Build complete! Binary: bin/$(BINARY_NAME)$(NC)"

## build-migrate: Build the migration tool
build-migrate:
	@echo "$(GREEN)Building migration tool...$(NC)"
	$(GOBUILD) -o bin/$(MIGRATE_BINARY) -v $(MIGRATE_PATH)
	@echo "$(GREEN)✓ Migration tool built: bin/$(MIGRATE_BINARY)$(NC)"

## run: Run the application
run: build
	@echo "$(GREEN)Starting application...$(NC)"
	./bin/$(BINARY_NAME)

## dev: Run with hot reload using air
dev:
	@echo "$(GREEN)Starting development server with hot reload...$(NC)"
	@which air > /dev/null || (echo "$(RED)Error: air not installed. Run 'make install-tools'$(NC)" && exit 1)
	air

## clean: Remove build artifacts and temporary files
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	$(GOCLEAN)
	rm -rf bin/
	rm -rf tmp/
	@echo "$(GREEN)✓ Clean complete!$(NC)"

## db-clean: Remove database file
db-clean:
	@echo "$(YELLOW)Removing database file...$(NC)"
	rm -f $(DATABASE_PATH)
	@echo "$(GREEN)✓ Database removed!$(NC)"

## migrate-up: Run database migrations up
migrate-up:
	@echo "$(GREEN)Running migrations up...$(NC)"
	$(GOBUILD) -o bin/$(MIGRATE_BINARY) $(MIGRATE_PATH)
	./bin/$(MIGRATE_BINARY) up
	@echo "$(GREEN)✓ Migrations complete!$(NC)"

## migrate-down: Run database migrations down
migrate-down:
	@echo "$(YELLOW)Running migrations down...$(NC)"
	$(GOBUILD) -o bin/$(MIGRATE_BINARY) $(MIGRATE_PATH)
	./bin/$(MIGRATE_BINARY) down
	@echo "$(GREEN)✓ Migrations rolled back!$(NC)"

## migrate-create: Create new migration files (usage: make migrate-create NAME=create_users_table)
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: NAME is required. Usage: make migrate-create NAME=create_users_table$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Creating migration: $(NAME)$(NC)"
	@timestamp=$$(date +%s); \
	up_file="$(MIGRATIONS_PATH)/$${timestamp}_$(NAME).up.sql"; \
	down_file="$(MIGRATIONS_PATH)/$${timestamp}_$(NAME).down.sql"; \
	touch $$up_file $$down_file; \
	echo "-- Migration: $(NAME)" > $$up_file; \
	echo "-- Add your SQL here" >> $$up_file; \
	echo "" >> $$up_file; \
	echo "-- Migration: $(NAME)" > $$down_file; \
	echo "-- Add your rollback SQL here" >> $$down_file; \
	echo "" >> $$down_file; \
	echo "$(GREEN)✓ Created migration files:$(NC)"; \
	echo "  $$up_file"; \
	echo "  $$down_file"

## db-reset: Reset database (down and up migrations)
db-reset: migrate-down migrate-up
	@echo "$(GREEN)✓ Database reset complete!$(NC)"

## test: Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GOTEST) -v ./...

## coverage: Run tests with coverage report
coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

## fmt: Format Go code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GOFMT) ./...
	@echo "$(GREEN)✓ Code formatted!$(NC)"

## vet: Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GOVET) ./...
	@echo "$(GREEN)✓ Vet complete!$(NC)"

## lint: Run golangci-lint
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)Error: golangci-lint not installed. Run 'make install-tools'$(NC)" && exit 1)
	golangci-lint run ./...
	@echo "$(GREEN)✓ Lint complete!$(NC)"

## swagger: Generate Swagger documentation
swagger:
	@echo "$(GREEN)Generating Swagger documentation...$(NC)"
	@which swag > /dev/null || (echo "$(RED)Error: swag not installed. Run 'make install-tools'$(NC)" && exit 1)
	swag init -g cmd/server/main.go -o docs
	@echo "$(GREEN)✓ Swagger docs generated!$(NC)"

## docker-build: Build Docker image
docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -f ci/docker/Dockerfile -t $(BINARY_NAME):latest .
	@echo "$(GREEN)✓ Docker image built!$(NC)"

## docker-run: Run application in Docker
docker-run:
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	docker-compose -f ci/docker/docker-compose.yml up -d
	@echo "$(GREEN)✓ Application running in Docker!$(NC)"

## docker-stop: Stop Docker containers
docker-stop:
	@echo "$(YELLOW)Stopping Docker containers...$(NC)"
	docker-compose -f ci/docker/docker-compose.yml down
	@echo "$(GREEN)✓ Containers stopped!$(NC)"

## check: Run all checks (fmt, vet, test)
check: fmt vet test
	@echo "$(GREEN)✓ All checks passed!$(NC)"

## setup: Initial project setup
setup: deps install-tools migrate-up
	@echo "$(GREEN)✓ Project setup complete!$(NC)"
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  1. Copy .env.example to .env and configure your settings"
	@echo "  2. Run 'make dev' to start the development server"
	@echo "  3. Visit http://localhost:8080 to access the API"

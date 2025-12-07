# Makefile for Go RESTful API Project
# This file contains common commands to manage your project

# Variables
BINARY_NAME=ginflow
MAIN_PATH=./cmd/server

# Version info (can be overridden during build)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-ldflags "-X github.com/alireza-akbarzadeh/ginflow/cmd/cli.Version=$(VERSION) -X github.com/alireza-akbarzadeh/ginflow/cmd/cli.GitCommit=$(GIT_COMMIT) -X github.com/alireza-akbarzadeh/ginflow/cmd/cli.BuildTime=$(BUILD_TIME)"

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

.PHONY: all build clean test coverage deps run dev docker-build docker-run help migrate reset-db version

# Default target
all: deps build

## help: Display this help message
help:
	@echo "$(GREEN)Available commands:$(NC)"
	@echo ""
	@echo "$(YELLOW)Setup & Dependencies:$(NC)"
	@echo "  make deps          - Download and tidy Go dependencies"
	@echo "  make install-tools - Install development tools (air, swag, golangci-lint)"
	@echo ""
	@echo "$(YELLOW)Database:$(NC)"
	@echo "  make migrate       - Run database migrations"
	@echo "  make reset-db      - Reset database (drop all tables and re-migrate)"
	@echo "  make drop-db       - Drop all database tables"
	@echo "  make db-status     - Check database connection status"
	@echo ""
	@echo "$(YELLOW)Development:$(NC)"
	@echo "  make run           - Run the application (ginflow serve)"
	@echo "  make dev           - Run with hot reload (requires air)"
	@echo "  make build         - Build the application binary"
	@echo "  make clean         - Remove build artifacts and temporary files"
	@echo "  make version       - Show version information"
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

## migrate: Run database migrations
migrate: build
	@echo "$(GREEN)Running database migrations...$(NC)"
	./bin/$(BINARY_NAME) migrate

## reset-db: Reset database (drop all and re-migrate)
reset-db: build
	@echo "$(YELLOW)Resetting database...$(NC)"
	./bin/$(BINARY_NAME) db reset --force

## drop-db: Drop all database tables
drop-db: build
	@echo "$(RED)Dropping all database tables...$(NC)"
	./bin/$(BINARY_NAME) db drop --force

## db-status: Check database connection status
db-status: build
	./bin/$(BINARY_NAME) db status

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
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@echo "$(GREEN)✓ All tools installed!$(NC)"

## build: Build the application binary
build:
	@echo "$(GREEN)Building application...$(NC)"
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) -v $(MAIN_PATH)
	@echo "$(GREEN)✓ Build complete! Binary: bin/$(BINARY_NAME)$(NC)"

## run: Run the application
run: build
	@echo "$(GREEN)Starting application...$(NC)"
	./bin/$(BINARY_NAME) serve

## version: Show version information
version: build
	./bin/$(BINARY_NAME) version

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
	docker-compose up -d
	@echo "$(GREEN)✓ Application running in Docker!$(NC)"

## docker-stop: Stop Docker containers
docker-stop:
	@echo "$(YELLOW)Stopping Docker containers...$(NC)"
	docker-compose down
	@echo "$(GREEN)✓ Containers stopped!$(NC)"

## check: Run all checks (fmt, vet, test)
check: fmt vet test
	@echo "$(GREEN)✓ All checks passed!$(NC)"

## setup: Initial project setup
setup: deps install-tools
	@echo "$(GREEN)✓ Project setup complete!$(NC)"
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  1. Copy .env.example to .env and configure your settings"
	@echo "  2. Run 'make dev' to start the development server"
	@echo "  3. Visit http://localhost:8080 to access the API"

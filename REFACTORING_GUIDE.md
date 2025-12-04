# Project Structure Refactoring Guide

## Overview

The project has been restructured to follow best practices for Go web applications with a clean separation of concerns.

## New Directory Structure

```
.
├── ci/
│   └── docker/                 # Docker configuration
│       ├── Dockerfile
│       └── docker-compose.yml
├── cmd/
│   ├── api/                    # Main application entry point
│   │   ├── main_new.go        # New refactored main file
│   │   └── main.go            # Old main file (to be replaced)
│   └── migrate/                # Database migration tool
├── pkg/                        # Public packages (reusable code)
│   ├── api/                    # API layer
│   │   ├── handlers/           # HTTP request handlers
│   │   │   ├── handler.go     # Base handler with dependencies
│   │   │   ├── auth.go        # Authentication handlers
│   │   │   ├── event.go       # Event handlers
│   │   │   └── attendee.go    # Attendee handlers
│   │   ├── middleware/         # HTTP middleware
│   │   │   ├── auth.go        # JWT authentication middleware
│   │   │   ├── logger.go      # Request logging middleware
│   │   │   └── cors.go        # CORS middleware
│   │   ├── routers/            # Route definitions
│   │   │   └── router.go      # Main router setup
│   │   └── helpers/            # Helper utilities
│   │       ├── context.go     # Context helpers
│   │       ├── response.go    # Response helpers
│   │       └── params.go      # Parameter parsing helpers
│   ├── repository/             # Data access layer
│   │   ├── repository.go      # Repository initialization
│   │   ├── user.go            # User repository
│   │   ├── event.go           # Event repository
│   │   └── attendee.go        # Attendee repository
│   └── config/                 # Configuration utilities
│       └── env.go             # Environment variable helpers
├── internal/                   # Private packages (deprecated)
│   └── ...                    # Old structure (to be removed)
└── docs/                       # Swagger documentation
```

## Key Changes

### 1. Eliminated `internal/` Directory

- Moved all code to `pkg/` for better reusability
- `pkg/` contains code that can be imported by other projects
- Better aligns with Go community practices

### 2. Organized API Layer (`pkg/api/`)

#### **Handlers** (`pkg/api/handlers/`)

- All HTTP request handlers organized by domain
- `handler.go`: Base handler struct with shared dependencies
- `auth.go`: Registration and login handlers
- `event.go`: CRUD operations for events
- `attendee.go`: Attendee management handlers

#### **Middleware** (`pkg/api/middleware/`)

- `auth.go`: JWT token verification and user authentication
- `logger.go`: Request/response logging
- `cors.go`: Cross-Origin Resource Sharing configuration

#### **Routers** (`pkg/api/routers/`)

- `router.go`: Centralized route definitions
- Groups routes by functionality (public vs protected)
- Applies middleware appropriately

#### **Helpers** (`pkg/api/helpers/`)

- `context.go`: User context management
- `response.go`: Standardized JSON responses
- `params.go`: Parameter parsing utilities

### 3. Repository Layer (`pkg/repository/`)

- Clean separation of data access logic
- Each entity has its own repository file
- Consistent naming: `*Repository` instead of `*Model`
- Better structured with clear method names
- Added helper methods for common operations

### 4. Configuration (`pkg/config/`)

- Centralized environment variable management
- Type-safe configuration helpers
- Easy to extend with new config options

### 5. Docker Files Moved to `ci/docker/`

- Better organization following CI/CD practices
- Dockerfile and docker-compose.yml in dedicated directory
- Makefile updated to reference correct paths

## Migration Steps

### Step 1: Backup Old Files

```bash
# Already done - old files remain in place
```

### Step 2: Replace main.go

```bash
mv cmd/api/main.go cmd/api/main_old.go
mv cmd/api/main_new.go cmd/api/main.go
```

### Step 3: Update Imports Throughout Project

The new import paths:

```go
// Old
"github.com/alireza-akbarzadeh/restful-app/internal/database"
"github.com/alireza-akbarzadeh/restful-app/internal/env"

// New
"github.com/alireza-akbarzadeh/restful-app/pkg/repository"
"github.com/alireza-akbarzadeh/restful-app/pkg/config"
"github.com/alireza-akbarzadeh/restful-app/pkg/api/handlers"
"github.com/alireza-akbarzadeh/restful-app/pkg/api/middleware"
"github.com/alireza-akbarzadeh/restful-app/pkg/api/routers"
"github.com/alireza-akbarzadeh/restful-app/pkg/api/helpers"
```

### Step 4: Build and Test

```bash
# Clean old build artifacts
make clean

# Download dependencies
go mod tidy

# Build the application
make build

# Run tests
make test

# Run the application
make run
```

### Step 5: Update Swagger Documentation

```bash
make swagger
```

### Step 6: Remove Old Files (After Testing)

```bash
# Once everything works
rm -rf internal/
rm cmd/api/main_old.go
rm cmd/api/auth.go
rm cmd/api/event.go
rm cmd/api/middleware.go
rm cmd/api/routes.go
rm cmd/api/server.go
rm cmd/api/context.go
```

## Benefits of New Structure

### 1. **Better Separation of Concerns**

- Each package has a single, well-defined responsibility
- Easier to understand and maintain
- Reduces coupling between components

### 2. **Improved Testability**

- Dependencies are explicitly injected
- Easy to mock repositories and services
- Handlers can be tested in isolation

### 3. **Scalability**

- Easy to add new features
- Clear patterns for where code belongs
- Supports microservices extraction if needed

### 4. **Code Reusability**

- `pkg/` code can be imported by other projects
- Helpers and utilities are easily accessible
- Repository pattern allows database swapping

### 5. **Industry Standard**

- Follows Go community best practices
- Similar to popular open-source Go projects
- Easier for new developers to understand

## Usage Examples

### Adding a New Handler

1. Create handler method in appropriate file:

```go
// pkg/api/handlers/event.go
func (h *Handler) SearchEvents(c *gin.Context) {
    query := c.Query("q")
    events, err := h.Repos.Events.Search(query)
    // ... handle response
}
```

2. Add route in router:

```go
// pkg/api/routers/router.go
events.GET("/search", handler.SearchEvents)
```

### Adding a New Repository Method

```go
// pkg/repository/event.go
func (r *EventRepository) Search(query string) ([]*Event, error) {
    // Implementation
}
```

### Adding a New Middleware

```go
// pkg/api/middleware/rate_limit.go
func RateLimit(requestsPerMinute int) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Implementation
        c.Next()
    }
}
```

## Testing the New Structure

```bash
# Test compilation
go build ./cmd/api

# Run specific package tests
go test ./pkg/repository/...
go test ./pkg/api/handlers/...
go test ./pkg/api/middleware/...

# Run all tests
make test

# Check test coverage
make coverage
```

## Troubleshooting

### Import Errors

Run `go mod tidy` to update dependencies.

### Build Errors

Ensure all old files are removed and new imports are used.

### Database Issues

The database schema hasn't changed, so existing databases work as-is.

### Swagger Documentation

Regenerate with `make swagger` after structural changes.

## Next Steps

1. **Add Unit Tests**: Write tests for all new handlers and repositories
2. **Add Integration Tests**: Test end-to-end flows
3. **Add Logging**: Implement structured logging throughout
4. **Add Metrics**: Add Prometheus metrics for monitoring
5. **Add Tracing**: Implement distributed tracing
6. **Add Validation**: Use a validation library for request validation
7. **Add Documentation**: Document all public APIs and methods

## Questions?

- Review the code in `pkg/` to understand the new structure
- Check `QUICKSTART.md` for development guidelines
- See `README.md` for general project information

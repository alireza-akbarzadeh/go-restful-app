# GinFlow - Modern Go API Framework

A modern RESTful API framework built with Go and Gin for managing events and attendees. Perfect for learning Go web development with best practices and clean architecture.

## 📋 Table of Contents



- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Available Commands](#available-commands)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [Best Practices](#best-practices)
- [Contributing](#contributing)

## ✨ Features

- ✅ RESTful API design
- ✅ JWT authentication
- ✅ SQLite database with migrations
- ✅ Swagger/OpenAPI documentation
- ✅ Hot reload development with Air
- ✅ Docker support
- ✅ Comprehensive Makefile for easy project management
- ✅ Clean architecture with separation of concerns
- ✅ Event CRUD operations
- ✅ User management and authentication
- ✅ Attendee management

## 📁 Project Structure

```
.
├── cmd/
│   ├── api/                    # Main API application
│   │   ├── main.go            # Application entry point
│   │   ├── server.go          # HTTP server setup
│   │   ├── routes.go          # Route definitions
│   │   ├── middleware.go      # Custom middleware
│   │   ├── auth.go            # Authentication handlers
│   │   ├── event.go           # Event handlers
│   │   └── context.go         # Context helpers
│   └── migrate/               # Database migration tool
│       ├── main.go            # Migration runner
│       └── migrations/        # SQL migration files
├── internal/
│   ├── database/              # Database layer
│   │   ├── models.go          # Data models
│   │   ├── users.go           # User repository
│   │   ├── events.go          # Event repository
│   │   └── attendees.go       # Attendee repository
│   ├── env/                   # Environment configuration
│   │   └── env.go             # Environment helpers
│   └── messages/              # Response messages
│       └── message.go         # Message templates
├── docs/                      # Swagger documentation (auto-generated)
├── .env.example               # Environment variables template
├── .gitignore                 # Git ignore rules
├── Makefile                   # Project automation commands
├── Dockerfile                 # Docker image definition
├── docker-compose.yml         # Docker Compose setup
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## 🔧 Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.21 or higher** - [Download Go](https://golang.org/dl/)
- **Make** - Usually pre-installed on macOS/Linux
- **Git** - [Download Git](https://git-scm.com/)
- **Docker** (optional) - [Download Docker](https://www.docker.com/)

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/alireza-akbarzadeh/ginflow.git
cd ginflow
```

### 2. Quick Setup

Run the automated setup command:

```bash
make setup
```

This will:

- Download all Go dependencies
- Install development tools (air, swag, migrate)
- Run database migrations

### 3. Configure Environment

Create your `.env` file from the example:

```bash
cp .env.example .env
```

Edit `.env` and update the values:

```env
PORT=8080
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
DATABASE_PATH=./data.db
APP_ENV=development
```

### 4. Start the Application

For development with hot reload:

```bash
make dev
```

Or run normally:

```bash
make run
```

The API will be available at `http://localhost:8080`

## 📝 Available Commands

Run `make help` to see all available commands. Here are the most common ones:

### Development Commands

```bash
make dev              # Start with hot reload
make run              # Build and run the application
make build            # Build the binary
make clean            # Clean build artifacts
```

### Database Commands

```bash
make migrate-up       # Run database migrations
make migrate-down     # Rollback migrations
make migrate-create NAME=create_something_table  # Create new migration
make db-reset         # Reset database (down + up)
make db-clean         # Remove database file
```

### Code Quality Commands

```bash
make test             # Run tests
make coverage         # Run tests with coverage report
make fmt              # Format code
make vet              # Run go vet
make lint             # Run linter
make swagger          # Generate Swagger docs
```

### Docker Commands

```bash
make docker-build     # Build Docker image
make docker-run       # Run in Docker
make docker-stop      # Stop Docker containers
```

## 📚 API Documentation

### Swagger UI

Once the application is running, access the interactive API documentation:

```
http://localhost:8080/swagger/index.html
```

### Main Endpoints

#### Authentication

```bash
# Register a new user
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}

# Login
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### Events (Requires Authentication)

```bash
# Create an event
POST /api/v1/events
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "name": "Tech Conference 2025",
  "description": "Annual tech conference",
  "location": "San Francisco",
  "dateTime": "2025-06-15T09:00:00Z"
}

# Get all events
GET /api/v1/events

# Get specific event
GET /api/v1/events/:id

# Update event
PUT /api/v1/events/:id
Authorization: Bearer <your-jwt-token>

# Delete event
DELETE /api/v1/events/:id
Authorization: Bearer <your-jwt-token>
```

#### Attendees

```bash
# Register for an event
POST /api/v1/events/:id/register
Authorization: Bearer <your-jwt-token>

# Cancel registration
DELETE /api/v1/events/:id/register
Authorization: Bearer <your-jwt-token>
```

## 🔨 Development

### Hot Reload Development

The project uses Air for hot reloading during development:

```bash
make dev
```

This will automatically rebuild and restart the server when you make changes to `.go` files.

### Creating a New Migration

When you need to change the database schema:

```bash
make migrate-create NAME=add_email_verification
```

This creates two files:

- `{timestamp}_add_email_verification.up.sql` - For applying changes
- `{timestamp}_add_email_verification.down.sql` - For rolling back

Edit these files with your SQL, then run:

```bash
make migrate-up
```

### Adding New Endpoints

1. **Define the handler** in `cmd/api/` (e.g., `event.go`)
2. **Add the route** in `cmd/api/routes.go`
3. **Update Swagger comments** above your handler function
4. **Regenerate Swagger docs**:
   ```bash
   make swagger
   ```

Example handler with Swagger documentation:

```go
// GetEvent retrieves a specific event by ID
// @Summary Get event by ID
// @Description Get details of a specific event
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} database.Event
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/events/{id} [get]
func (app *application) GetEvent(c *gin.Context) {
    // Your implementation here
}
```

## 🧪 Testing

### Run All Tests

```bash
make test
```

### Run Tests with Coverage

```bash
make coverage
```

This generates a `coverage.html` file you can open in your browser.

### Writing Tests

Create test files alongside your code with `_test.go` suffix:

```go
// event_test.go
package main

import (
    "testing"
)

func TestCreateEvent(t *testing.T) {
    // Your test implementation
}
```

## 🐳 Docker Deployment

### Build Docker Image

```bash
make docker-build
```

### Run with Docker Compose

```bash
make docker-run
```

### Stop Docker Containers

```bash
make docker-stop
```

## 🎯 Best Practices Implemented

### 1. **Project Organization**

- Clean separation of concerns
- `cmd/` for application entry points
- `internal/` for private application code
- `docs/` for auto-generated documentation

### 2. **Error Handling**

- Consistent error responses
- Proper HTTP status codes
- Meaningful error messages

### 3. **Security**

- JWT authentication
- Password hashing (bcrypt)
- Environment-based configuration
- No sensitive data in version control

### 4. **Database**

- Migration-based schema management
- Repository pattern for data access
- Connection pooling

### 5. **API Design**

- RESTful conventions
- Versioned endpoints (`/api/v1/`)
- Proper HTTP methods
- JSON request/response

### 6. **Development Workflow**

- Hot reload for faster development
- Automated tasks with Makefile
- Docker support for consistent environments
- Swagger documentation

## 🔍 Troubleshooting

### Port Already in Use

If port 8080 is already in use, update the `PORT` in your `.env` file:

```env
PORT=3000
```

### Migration Errors

Reset the database:

```bash
make db-reset
```

### Dependencies Issues

Clean and reinstall dependencies:

```bash
make clean
make deps
```

## 📖 Learning Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [Swagger/OpenAPI](https://swagger.io/)
- [Database Migration Best Practices](https://www.prisma.io/dataguide/types/relational/what-are-database-migrations)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 💡 Next Steps for Learning

1. **Add More Features**:

   - Email notifications
   - File uploads
   - Pagination for list endpoints
   - Search and filtering
   - Rate limiting

2. **Improve Testing**:

   - Add unit tests for handlers
   - Integration tests
   - Test coverage > 80%

3. **Add More Middleware**:

   - Request logging
   - CORS configuration
   - Rate limiting
   - Request ID tracking

4. **Database Enhancements**:

   - Add indexes for performance
   - Implement soft deletes
   - Add database seeding

5. **Production Ready**:
   - Add health check endpoint
   - Implement graceful shutdown
   - Add metrics and monitoring
   - Set up CI/CD pipeline

## 📞 Support

If you have questions or need help, please:

- Open an issue on GitHub
- Check existing issues for solutions
- Review the documentation

---

**Happy Coding! 🚀**

Made with ❤️ by Alireza Akbarzadeh

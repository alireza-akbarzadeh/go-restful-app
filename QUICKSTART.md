# Go Restful API - Quick Start Guide

Welcome! This guide will help you get started with your Go REST API project step by step.

## Step 1: Initial Setup (First Time Only)

### 1.1 Install Prerequisites

Make sure you have Go installed:

```bash
go version  # Should show Go 1.21 or higher
```

### 1.2 Setup the Project

Run the automated setup:

```bash
make setup
```

This command will:

- ✅ Download all Go dependencies
- ✅ Install development tools (air, swag, golangci-lint)

### 1.3 Configure Environment Variables

Copy the example environment file:

```bash
cp .env.example .env
```

Open `.env` and customize if needed:

```env
PORT=8080
JWT_SECRET=change-this-to-something-secure

# Database Configuration (PostgreSQL)
DATABASE_URL=postgres://username:password@localhost:5432/database_name
# Or for development with Docker:
# DATABASE_URL=postgres://user:password@db:5432/ginflow_dev
```

## Step 2: Daily Development Workflow

### Starting Development Server

Start the server with hot reload (automatically restarts when you change code):

```bash
make dev
```

Your API will be available at: `http://localhost:8080`

### Testing the API

#### Using the Swagger UI (Recommended for Beginners)

Open your browser and visit:

```
http://localhost:8080/swagger/index.html
```

Here you can:

- See all available endpoints
- Try making requests directly from the browser
- See request/response examples

#### Using curl or httpie

Register a new user:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

Login:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

Save the token from the response and use it for authenticated requests:

```bash
TOKEN="your-jwt-token-here"

curl -X GET http://localhost:8080/api/v1/events \
  -H "Authorization: Bearer $TOKEN"
```

## Step 3: Adding New Features

### 3.1 Creating a New Endpoint

Let's say you want to add a "Get User Profile" endpoint.

**Step 3.1.1: Add the handler function** in `pkg/api/handlers/` (create a new file like `profile.go`):

```go
package handlers

import (
    "net/http"
    "github.com/alireza-akbarzadeh/ginflow/pkg/api/helpers"
    "github.com/gin-gonic/gin"
)

// GetUserProfile returns the current user's profile
// @Summary Get user profile
// @Description Get the profile of the currently authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Profile
// @Failure 401 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Router /api/v1/profile [get]
func (h *Handler) GetUserProfile(c *gin.Context) {
    // Get authenticated user
    authUser := helpers.GetUserFromContext(c)
    if authUser == nil {
        helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
        return
    }

    // Get profile
    profile, err := h.Repos.Profiles.GetByUserID(authUser.ID)
    if err != nil {
        helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve profile")
        return
    }

    c.JSON(http.StatusOK, profile)
}
```

**Step 3.1.2: Add the route** in `pkg/api/routers/router.go`:

```go
// In the protected routes section
protected.GET("/profile", handler.GetUserProfile)
```

**Step 3.1.3: Update Swagger docs**:

```bash
make swagger
```

**Step 3.1.4: Test your new endpoint**:

Restart the server (if using `make dev`, it should auto-restart). Visit Swagger UI at `http://localhost:8080/swagger/index.html` to test!

### 3.2 Adding a Database Table

**Step 3.2.1: Create the model** in `pkg/models/` (create a new file like `comment.go`):

```go
package models

import "time"

type Comment struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    EventID   uint      `json:"eventId" gorm:"not null"`
    UserID    int       `json:"userId" gorm:"not null"`
    Content   string    `json:"content" gorm:"type:text;not null"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`

    // Relationships
    Event Event `json:"event,omitempty" gorm:"foreignKey:EventID"`
    User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
```

**Step 3.2.2: Create repository methods** in `pkg/repository/` (create a new file like `comment.go`):

```go
package repository

import (
    "github.com/alireza-akbarzadeh/ginflow/pkg/models"
    "gorm.io/gorm"
)

type CommentRepository struct {
    DB *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
    return &CommentRepository{DB: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
    return r.DB.Create(comment).Error
}

func (r *CommentRepository) GetByEventID(eventID uint) ([]*models.Comment, error) {
    var comments []*models.Comment
    err := r.DB.Preload("User").Where("event_id = ?", eventID).Order("created_at DESC").Find(&comments).Error
    return comments, err
}
```

**Step 3.2.3: Update Models struct** in `pkg/repository/repository.go`:

```go
type Models struct {
    Users     *UserRepository
    Events    *EventRepository
    Attendees *AttendeeRepository
    Categories *CategoryRepository
    Comments  *CommentRepository  // Add this
}

func NewModels(db *gorm.DB) *Models {
    return &Models{
        Users:      NewUserRepository(db),
        Events:     NewEventRepository(db),
        Attendees:  NewAttendeeRepository(db),
        Categories: NewCategoryRepository(db),
        Comments:   NewCommentRepository(db),  // Add this
    }
}
```

**Step 3.2.4: Update AutoMigrate** in `cmd/server/main.go`:

```go
err := db.AutoMigrate(
    &models.User{},
    &models.Event{},
    &models.Attendee{},
    &models.Category{},
    &models.Comment{},  // Add this
)
```

GORM will automatically create the table and relationships when you restart the server!

## Step 4: Common Tasks

### Run Tests

```bash
make test
```

### Check Code Quality

```bash
make fmt   # Format code
make vet   # Check for issues
make lint  # Run linter
```

### Reset Database

If you mess up the database, restart the server - GORM AutoMigrate will recreate tables automatically.

### Clean Build Artifacts

```bash
make clean
```

### See All Available Commands

```bash
make help
```

## Step 5: Debugging Tips

### Check if Server is Running

```bash
curl http://localhost:8080/health
```

### View Logs

The server logs appear in your terminal where you ran `make dev`

### Database Issues

Connect to your PostgreSQL database using your preferred database client (pgAdmin, DBeaver, etc.) or command line:

```bash
psql "your-connection-string"
```

Or use database management tools like:

- [pgAdmin](https://www.pgadmin.org/)
- [DBeaver](https://dbeaver.io/)
- [TablePlus](https://tableplus.com/)

### Port Already in Use

Change the port in `.env`:

```env
PORT=3000
```

## Step 6: Going to Production

### Build for Production

```bash
make build
```

This creates a binary in `bin/api-server`

### Using Docker

```bash
make docker-build
make docker-run
```

## Learning Path Recommendations

1. **Week 1-2**: Get comfortable with the existing CRUD operations

   - Understand how routes work in `pkg/api/routers/router.go`
   - Study the authentication flow in `pkg/api/handlers/auth.go`
   - Practice making API requests using Swagger UI

2. **Week 3-4**: Add your first feature

   - Create a simple new endpoint following the profile example
   - Add validation using the existing helpers
   - Write tests in the `tests/` directory

3. **Week 5-6**: Work with databases

   - Create a new model in `pkg/models/`
   - Add repository methods in `pkg/repository/`
   - Update AutoMigrate in `cmd/server/main.go`

4. **Week 7-8**: Advanced topics
   - Add pagination to list endpoints
   - Implement search and filtering
   - Add file uploads
   - Rate limiting and security

## Need Help?

- Check the main `README.md` for detailed documentation
- Look at existing code in `pkg/api/handlers/`, `pkg/models/`, and `pkg/repository/`
- Study the profile implementation as a complete example
- Search for Go/GORM tutorials for specific topics
- Ask questions in Go community forums

## Useful Commands Cheat Sheet

```bash
# Start development
make dev

# Run tests
make test

# Format code
make fmt

# See all commands
make help
```

**Happy coding! 🚀**

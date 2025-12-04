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
- ✅ Install development tools (air, swag, migrate, golangci-lint)
- ✅ Create the database and run migrations

### 1.3 Configure Environment Variables

Copy the example environment file:

```bash
cp .env.example .env
```

Open `.env` and customize if needed:

```env
PORT=8080
JWT_SECRET=change-this-to-something-secure
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

**Step 3.1.1: Add the handler function** in `cmd/api/auth.go` (or create a new file like `user.go`):

```go
// GetUserProfile returns the current user's profile
// @Summary Get user profile
// @Description Get the profile of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} database.User
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/user/profile [get]
func (app *application) GetUserProfile(c *gin.Context) {
    userID := c.GetInt64("userID") // From JWT middleware

    user, err := app.models.Users.GetById(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}
```

**Step 3.1.2: Add the route** in `cmd/api/routes.go`:

```go
// In the protected routes section
protected.GET("/user/profile", app.GetUserProfile)
```

**Step 3.1.3: Update Swagger docs**:

```bash
make swagger
```

**Step 3.1.4: Test your new endpoint**:

Restart the server (if using `make dev`, it should auto-restart). Visit Swagger UI to test!

### 3.2 Adding a Database Table

**Step 3.2.1: Create a migration**:

```bash
make migrate-create NAME=create_comments_table
```

This creates two files in `cmd/migrate/migrations/`:

- `{timestamp}_create_comments_table.up.sql` - for creating the table
- `{timestamp}_create_comments_table.down.sql` - for removing the table

**Step 3.2.2: Edit the `.up.sql` file**:

```sql
-- Migration: create_comments_table
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_comments_event_id ON comments(event_id);
```

**Step 3.2.3: Edit the `.down.sql` file**:

```sql
-- Migration: create_comments_table
DROP INDEX IF EXISTS idx_comments_event_id;
DROP TABLE IF EXISTS comments;
```

**Step 3.2.4: Run the migration**:

```bash
make migrate-up
```

**Step 3.2.5: Create the model** in `internal/database/models.go`:

```go
type Comment struct {
    ID        int64     `json:"id"`
    EventID   int64     `json:"event_id"`
    UserID    int64     `json:"user_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}
```

**Step 3.2.6: Create repository methods** in `internal/database/comments.go`:

```go
package database

import "database/sql"

type CommentModel struct {
    DB *sql.DB
}

func (m *CommentModel) Insert(comment *Comment) error {
    query := `INSERT INTO comments (event_id, user_id, content)
              VALUES (?, ?, ?)`

    result, err := m.DB.Exec(query, comment.EventID, comment.UserID, comment.Content)
    if err != nil {
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return err
    }

    comment.ID = id
    return nil
}

func (m *CommentModel) GetByEventID(eventID int64) ([]*Comment, error) {
    query := `SELECT id, event_id, user_id, content, created_at
              FROM comments WHERE event_id = ? ORDER BY created_at DESC`

    rows, err := m.DB.Query(query, eventID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []*Comment
    for rows.Next() {
        var c Comment
        err := rows.Scan(&c.ID, &c.EventID, &c.UserID, &c.Content, &c.CreatedAt)
        if err != nil {
            return nil, err
        }
        comments = append(comments, &c)
    }

    return comments, nil
}
```

**Step 3.2.7: Update Models struct** in `internal/database/models.go`:

```go
type Models struct {
    Users     UserModel
    Events    EventModel
    Attendees AttendeeModel
    Comments  CommentModel  // Add this
}

func NewModels(db *sql.DB) *Models {
    return &Models{
        Users:     UserModel{DB: db},
        Events:    EventModel{DB: db},
        Attendees: AttendeeModel{DB: db},
        Comments:  CommentModel{DB: db},  // Add this
    }
}
```

Now you can use `app.models.Comments` in your handlers!

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

If you mess up the database:

```bash
make db-reset
```

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

View the database:

```bash
sqlite3 data.db
sqlite> .tables
sqlite> SELECT * FROM users;
sqlite> .quit
```

Or use a GUI tool like [DB Browser for SQLite](https://sqlitebrowser.org/)

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

   - Understand how routes work
   - Study the authentication flow
   - Practice making API requests

2. **Week 3-4**: Add your first feature

   - Create a simple new endpoint
   - Add validation
   - Write tests

3. **Week 5-6**: Work with databases

   - Create a new table
   - Add relationships
   - Learn about indexes

4. **Week 7-8**: Advanced topics
   - Add pagination
   - Implement search
   - Add file uploads
   - Rate limiting

## Need Help?

- Check the main `README.md` for detailed documentation
- Look at existing code for examples
- Search for Go tutorials on specific topics
- Ask questions in Go community forums

## Useful Commands Cheat Sheet

```bash
# Start development
make dev

# Run migrations
make migrate-up

# Create new migration
make migrate-create NAME=my_migration

# Reset database
make db-reset

# Run tests
make test

# Format code
make fmt

# See all commands
make help
```

**Happy coding! 🚀**

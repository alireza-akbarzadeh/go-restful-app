# API Development Guide

This guide provides step-by-step instructions for creating new APIs in the GinFlow project. Follow these patterns to maintain consistency and code quality.

---

## Table of Contents

1. [Development Order](#development-order)
2. [Step 1: Define the Model](#step-1-define-the-model)
3. [Step 2: Create Repository Interface](#step-2-create-repository-interface)
4. [Step 3: Implement Repository](#step-3-implement-repository)
5. [Step 4: Register Repository](#step-4-register-repository)
6. [Step 5: Create Handler](#step-5-create-handler)
7. [Step 6: Define Routes](#step-6-define-routes)
8. [Step 7: Create Mock for Testing](#step-7-create-mock-for-testing)
9. [Step 8: Write Tests](#step-8-write-tests)
10. [Error Handling](#error-handling)
11. [Pagination Guide](#pagination-guide)
12. [Validation Guide](#validation-guide)
13. [Best Practices](#best-practices)
14. [Common Mistakes to Avoid](#common-mistakes-to-avoid)

---

## Development Order

Always follow this order when creating a new API resource:

```
1. Model          → Define data structure
2. Interface      → Define repository contract
3. Repository     → Implement database operations
4. Register       → Add to repository container
5. Handler        → Implement HTTP handlers
6. Routes         → Define API endpoints
7. Mock           → Create test mock
8. Tests          → Write unit tests
```

---

## Step 1: Define the Model

Location: `internal/models/`

### Example: Creating a `Book` model

```go
// internal/models/book.go
package models

import "time"

type Book struct {
    ID          int       `json:"id" gorm:"primaryKey"`
    UserID      int       `json:"user_id" gorm:"not null;index"`
    Title       string    `json:"title" gorm:"size:255;not null" binding:"required,min=2,max=255"`
    Author      string    `json:"author" gorm:"size:255;not null" binding:"required,min=2,max=255"`
    ISBN        string    `json:"isbn" gorm:"size:20;uniqueIndex" binding:"omitempty,len=13"`
    Description string    `json:"description" gorm:"type:text"`
    Price       float64   `json:"price" gorm:"type:decimal(10,2)" binding:"omitempty,gte=0"`
    PublishedAt *time.Time `json:"published_at"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    // Relationships
    User       User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Categories []Category `json:"categories,omitempty" gorm:"many2many:book_categories;"`
}

// TableName specifies the table name (optional, GORM auto-pluralizes)
func (Book) TableName() string {
    return "books"
}
```

### Model Guidelines

| Tag       | Purpose                | Example                            |
| --------- | ---------------------- | ---------------------------------- |
| `json`    | JSON field name        | `json:"title"`                     |
| `gorm`    | Database configuration | `gorm:"size:255;not null;index"`   |
| `binding` | Validation rules       | `binding:"required,min=2,max=255"` |

### Common GORM Tags

```go
gorm:"primaryKey"           // Primary key
gorm:"not null"             // NOT NULL constraint
gorm:"uniqueIndex"          // Unique index
gorm:"index"                // Regular index
gorm:"size:255"             // VARCHAR size
gorm:"type:text"            // Text type
gorm:"type:decimal(10,2)"   // Decimal type
gorm:"default:0"            // Default value
gorm:"foreignKey:UserID"    // Foreign key relation
gorm:"many2many:table_name" // Many-to-many relation
```

### Common Validation Tags

```go
binding:"required"          // Field is required
binding:"min=2"             // Minimum length/value
binding:"max=255"           // Maximum length/value
binding:"email"             // Must be valid email
binding:"url"               // Must be valid URL
binding:"gte=0"             // Greater than or equal
binding:"lte=100"           // Less than or equal
binding:"oneof=active inactive" // Must be one of values
binding:"omitempty"         // Skip validation if empty
```

---

## Step 2: Create Repository Interface

Location: `internal/repository/interfaces/`

```go
// internal/repository/interfaces/book.go
package interfaces

import (
    "context"

    "github.com/alireza-akbarzadeh/ginflow/internal/models"
    "github.com/alireza-akbarzadeh/ginflow/internal/query"
)

type BookRepositoryInterface interface {
    // Create
    Insert(ctx context.Context, book *models.Book) (*models.Book, error)

    // Read
    Get(ctx context.Context, id int) (*models.Book, error)
    GetByISBN(ctx context.Context, isbn string) (*models.Book, error)
    GetAll(ctx context.Context, req *query.QueryParams) ([]models.Book, *query.PaginatedList, error)
    GetByUser(ctx context.Context, userID int) ([]models.Book, error)

    // Update
    Update(ctx context.Context, book *models.Book) error

    // Delete
    Delete(ctx context.Context, id int) error
}
```

### Interface Guidelines

- Always use `context.Context` as first parameter
- Return pointers for single objects (`*models.Book`)
- Return slices for collections (`[]models.Book`)
- Return `error` as last return value
- Include pagination method for list operations

---

## Step 3: Implement Repository

Location: `internal/repository/`

```go
// internal/repository/book.go
package repository

import (
    "context"

    "github.com/alireza-akbarzadeh/ginflow/internal/models"
    "github.com/alireza-akbarzadeh/ginflow/internal/query"
    "gorm.io/gorm"
)

type BookRepository struct {
    DB *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
    return &BookRepository{DB: db}
}

// Insert creates a new book
func (r *BookRepository) Insert(ctx context.Context, book *models.Book) (*models.Book, error) {
    result := r.DB.WithContext(ctx).Create(book)
    if result.Error != nil {
        return nil, result.Error
    }
    return book, nil
}

// Get retrieves a book by ID
func (r *BookRepository) Get(ctx context.Context, id int) (*models.Book, error) {
    var book models.Book
    result := r.DB.WithContext(ctx).
        Preload("User").
        Preload("Categories").
        First(&book, id)

    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil // Return nil, nil for not found (not an error)
        }
        return nil, result.Error
    }
    return &book, nil
}

// GetByISBN retrieves a book by ISBN
func (r *BookRepository) GetByISBN(ctx context.Context, isbn string) (*models.Book, error) {
    var book models.Book
    result := r.DB.WithContext(ctx).Where("isbn = ?", isbn).First(&book)

    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, result.Error
    }
    return &book, nil
}

// GetAll retrieves books with advanced filtering, sorting, and pagination
func (r *BookRepository) GetAll(ctx context.Context, req *query.QueryParams) ([]models.Book, *query.PaginatedList, error) {
    var books []models.Book
    var total int64

    // Build query with security controls
    builder := query.NewQueryBuilder(r.DB.WithContext(ctx).Model(&models.Book{})).
        WithRequest(req).
        AllowFilters("title", "author", "user_id", "price", "created_at").  // Whitelist filter fields
        AllowSorts("title", "author", "price", "created_at", "updated_at"). // Whitelist sort fields
        SearchColumns("title", "author", "description").                     // Searchable columns
        DefaultSort("created_at", query.SortDesc)                            // Default sort order

    // Get count if needed
    if req.IncludeTotal {
        countQuery := r.DB.WithContext(ctx).Model(&models.Book{})
        for _, filter := range req.Filters {
            countQuery = query.FilterBy(filter)(countQuery)
        }
        if req.Search != "" {
            countQuery = query.Search(req.Search, "title", "author", "description")(countQuery)
        }
        countQuery.Count(&total)
    }

    // Execute query (use dbQuery to avoid shadowing package name)
    dbQuery := builder.Build()
    if err := dbQuery.Preload("User").Preload("Categories").Find(&books).Error; err != nil {
        return nil, nil, err
    }

    // Get cursor IDs
    var firstID, lastID int
    if len(books) > 0 {
        firstID = books[0].ID
        lastID = books[len(books)-1].ID
    }

    // Build response with HATEOAS links
    result := query.BuildResponse(books, req, total, len(books), firstID, lastID)

    return books, result, nil
}

// Update updates an existing book
func (r *BookRepository) Update(ctx context.Context, book *models.Book) error {
    return r.DB.WithContext(ctx).Save(book).Error
}

// Delete removes a book by ID
func (r *BookRepository) Delete(ctx context.Context, id int) error {
    return r.DB.WithContext(ctx).Delete(&models.Book{}, id).Error
}

// GetByUser retrieves all books by a user
func (r *BookRepository) GetByUser(ctx context.Context, userID int) ([]models.Book, error) {
    var books []models.Book
    result := r.DB.WithContext(ctx).
        Where("user_id = ?", userID).
        Preload("Categories").
        Find(&books)

    if result.Error != nil {
        return nil, result.Error
    }
    return books, nil
}
```

### Repository Guidelines

- **Always use `WithContext(ctx)`** for all database operations
- **Return `nil, nil`** when record not found (not an error condition)
- **Use `Preload()`** for eager loading relationships
- **Whitelist filter/sort fields** to prevent SQL injection

---

## Step 4: Register Repository

Location: `internal/repository/repository.go`

```go
// Add to the Repositories struct
type Repositories struct {
    Users      interfaces.UserRepositoryInterface
    Events     interfaces.EventRepositoryInterface
    Products   interfaces.ProductRepositoryInterface
    Categories interfaces.CategoryRepositoryInterface
    Comments   interfaces.CommentRepositoryInterface
    Profiles   interfaces.ProfileRepositoryInterface
    Attendees  interfaces.AttendeeRepositoryInterface
    Baskets    interfaces.BasketRepositoryInterface
    Books      interfaces.BookRepositoryInterface  // Add this line
}

// Add to NewRepositories function
func NewRepositories(db *gorm.DB) *Repositories {
    return &Repositories{
        Users:      NewUserRepository(db),
        Events:     NewEventRepository(db),
        Products:   NewProductRepository(db),
        Categories: NewCategoryRepository(db),
        Comments:   NewCommentRepository(db),
        Profiles:   NewProfileRepository(db),
        Attendees:  NewAttendeeRepository(db),
        Baskets:    NewBasketRepository(db),
        Books:      NewBookRepository(db),  // Add this line
    }
}
```

---

## Step 5: Create Handler

Location: `internal/api/handlers/`

```go
// internal/api/handlers/book.go
package handlers

import (
    "net/http"

    "github.com/alireza-akbarzadeh/ginflow/internal/api/helpers"
    appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
    "github.com/alireza-akbarzadeh/ginflow/internal/logging"
    "github.com/alireza-akbarzadeh/ginflow/internal/models"
    "github.com/alireza-akbarzadeh/ginflow/internal/query"
    "github.com/gin-gonic/gin"
)

// CreateBook handles book creation
// @Summary      Create a new book
// @Description  Create a new book (requires authentication)
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        book  body      models.Book  true  "Book object"
// @Success      201   {object}  models.Book
// @Failure      400   {object}  helpers.ErrorResponse
// @Failure      401   {object}  helpers.ErrorResponse
// @Failure      500   {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/books [post]
func (h *Handler) CreateBook(c *gin.Context) {
    ctx := c.Request.Context()

    // Bind and validate JSON
    var book models.Book
    if !helpers.BindJSON(c, &book) {
        return // BindJSON already sent error response
    }

    // Get authenticated user
    user, ok := helpers.GetAuthenticatedUser(c)
    if !ok {
        return // Already sent 401 response
    }
    book.UserID = user.ID

    // Log the operation
    logging.Debug(ctx, "creating new book", "title", book.Title, "user_id", user.ID)

    // Check for duplicate ISBN
    if book.ISBN != "" {
        existing, err := h.Repos.Books.GetByISBN(ctx, book.ISBN)
        if helpers.HandleError(c, err, "Failed to check ISBN") {
            return
        }
        if existing != nil {
            helpers.RespondWithAppError(c, appErrors.Newf(appErrors.ErrAlreadyExists, "book with ISBN '%s' already exists", book.ISBN), "")
            return
        }
    }

    // Create the book
    createdBook, err := h.Repos.Books.Insert(ctx, &book)
    if helpers.HandleError(c, err, "Failed to create book") {
        return
    }

    logging.Info(ctx, "book created successfully", "book_id", createdBook.ID, "title", createdBook.Title)
    c.JSON(http.StatusCreated, createdBook)
}

// GetAllBooks retrieves all books with advanced pagination
// @Summary      Get all books
// @Description  Get books with filtering, sorting, and pagination
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        page        query     int     false  "Page number (default: 1)"
// @Param        page_size   query     int     false  "Page size (default: 20, max: 100)"
// @Param        sort        query     string  false  "Sort: '-created_at,title:asc'"
// @Param        search      query     string  false  "Search in title, author, description"
// @Param        title[like] query     string  false  "Filter by title (partial match)"
// @Param        author[eq]  query     string  false  "Filter by exact author"
// @Param        price[gte]  query     number  false  "Filter by minimum price"
// @Param        price[lte]  query     number  false  "Filter by maximum price"
// @Success      200         {object}  query.PaginatedList{data=[]models.Book}
// @Failure      500         {object}  helpers.ErrorResponse
// @Router       /api/v1/books [get]
func (h *Handler) GetAllBooks(c *gin.Context) {
    ctx := c.Request.Context()

    // Parse query parameters (pagination, filtering, sorting, search)
    req := query.ParseFromContext(c)

    logging.Debug(ctx, "retrieving books",
        "page", req.Page,
        "page_size", req.PageSize,
        "search", req.Search,
    )

    books, result, err := h.Repos.Books.GetAll(ctx, req)
    if helpers.HandleError(c, err, "Failed to retrieve books") {
        return
    }

    logging.Debug(ctx, "books retrieved", "count", len(books))
    c.JSON(http.StatusOK, result)
}

// GetBook retrieves a single book by ID
// @Summary      Get a book
// @Description  Get a book by ID
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  models.Book
// @Failure      404  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Router       /api/v1/books/{id} [get]
func (h *Handler) GetBook(c *gin.Context) {
    ctx := c.Request.Context()

    // Parse ID parameter
    id, err := helpers.ParseIDParam(c, "id")
    if err != nil {
        return // Already sent error response
    }

    book, err := h.Repos.Books.Get(ctx, id)
    if helpers.HandleError(c, err, "Failed to retrieve book") {
        return
    }

    // Check if book exists (nil means not found)
    if book == nil {
        helpers.RespondWithAppError(c, appErrors.Newf(appErrors.ErrNotFound, "book with ID %d not found", id), "")
        return
    }

    c.JSON(http.StatusOK, book)
}

// UpdateBook updates a book
// @Summary      Update a book
// @Description  Update a book by ID (owner only)
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        id    path      int          true  "Book ID"
// @Param        book  body      models.Book  true  "Book object"
// @Success      200   {object}  models.Book
// @Failure      400   {object}  helpers.ErrorResponse
// @Failure      401   {object}  helpers.ErrorResponse
// @Failure      403   {object}  helpers.ErrorResponse
// @Failure      404   {object}  helpers.ErrorResponse
// @Failure      500   {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/books/{id} [put]
func (h *Handler) UpdateBook(c *gin.Context) {
    ctx := c.Request.Context()

    // Parse ID
    id, err := helpers.ParseIDParam(c, "id")
    if err != nil {
        return
    }

    // Get authenticated user
    user, ok := helpers.GetAuthenticatedUser(c)
    if !ok {
        return
    }

    logging.Debug(ctx, "updating book", "book_id", id, "user_id", user.ID)

    // Get existing book
    existingBook, err := h.Repos.Books.Get(ctx, id)
    if helpers.HandleError(c, err, "Failed to retrieve book") {
        return
    }
    if existingBook == nil {
        helpers.RespondWithAppError(c, appErrors.Newf(appErrors.ErrNotFound, "book with ID %d not found", id), "")
        return
    }

    // Check ownership
    if existingBook.UserID != user.ID {
        helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrForbidden, "you do not have permission to update this book"), "")
        return
    }

    // Bind update data
    var updateData models.Book
    if err := c.ShouldBindJSON(&updateData); err != nil {
        helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
        return
    }

    // Update fields (preserve ID and UserID)
    existingBook.Title = updateData.Title
    existingBook.Author = updateData.Author
    existingBook.ISBN = updateData.ISBN
    existingBook.Description = updateData.Description
    existingBook.Price = updateData.Price
    existingBook.PublishedAt = updateData.PublishedAt

    if err := h.Repos.Books.Update(ctx, existingBook); err != nil {
        helpers.HandleError(c, err, "Failed to update book")
        return
    }

    logging.Info(ctx, "book updated successfully", "book_id", id)
    c.JSON(http.StatusOK, existingBook)
}

// DeleteBook deletes a book
// @Summary      Delete a book
// @Description  Delete a book by ID (owner only)
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      204  {object}  nil
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      403  {object}  helpers.ErrorResponse
// @Failure      404  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/books/{id} [delete]
func (h *Handler) DeleteBook(c *gin.Context) {
    ctx := c.Request.Context()

    id, err := helpers.ParseIDParam(c, "id")
    if err != nil {
        return
    }

    user, ok := helpers.GetAuthenticatedUser(c)
    if !ok {
        return
    }

    logging.Debug(ctx, "deleting book", "book_id", id, "user_id", user.ID)

    // Get existing book
    existingBook, err := h.Repos.Books.Get(ctx, id)
    if helpers.HandleError(c, err, "Failed to retrieve book") {
        return
    }
    if existingBook == nil {
        helpers.RespondWithAppError(c, appErrors.Newf(appErrors.ErrNotFound, "book with ID %d not found", id), "")
        return
    }

    // Check ownership
    if existingBook.UserID != user.ID {
        helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrForbidden, "you do not have permission to delete this book"), "")
        return
    }

    if err := h.Repos.Books.Delete(ctx, id); err != nil {
        helpers.HandleError(c, err, "Failed to delete book")
        return
    }

    logging.Info(ctx, "book deleted successfully", "book_id", id)
    c.Status(http.StatusNoContent)
}
```

### Handler Pattern Summary

```go
func (h *Handler) HandlerName(c *gin.Context) {
    ctx := c.Request.Context()

    // 1. Parse parameters (ID, query params)
    id, err := helpers.ParseIDParam(c, "id")
    if err != nil {
        return
    }

    // 2. Get authenticated user (if required)
    user, ok := helpers.GetAuthenticatedUser(c)
    if !ok {
        return
    }

    // 3. Bind JSON body (if required)
    var data models.Model
    if !helpers.BindJSON(c, &data) {
        return
    }

    // 4. Log the operation
    logging.Debug(ctx, "operation description", "key", value)

    // 5. Business logic (check existence, ownership, etc.)
    existing, err := h.Repos.Resource.Get(ctx, id)
    if helpers.HandleError(c, err, "Failed to retrieve resource") {
        return
    }
    if existing == nil {
        helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrNotFound, "not found"), "")
        return
    }

    // 6. Perform operation
    result, err := h.Repos.Resource.Insert(ctx, &data)
    if helpers.HandleError(c, err, "Failed to create resource") {
        return
    }

    // 7. Log success and respond
    logging.Info(ctx, "operation successful", "id", result.ID)
    c.JSON(http.StatusOK, result)
}
```

---

## Step 6: Define Routes

Location: `internal/api/routers/`

### Create router file

```go
// internal/api/routers/book.go
package routers

import (
    "github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
    "github.com/alireza-akbarzadeh/ginflow/internal/api/middleware"
    "github.com/gin-gonic/gin"
)

func RegisterBookRoutes(rg *gin.RouterGroup, h *handlers.Handler, authMiddleware *middleware.AuthMiddleware) {
    books := rg.Group("/books")
    {
        // Public routes (no auth required)
        books.GET("", h.GetAllBooks)
        books.GET("/:id", h.GetBook)

        // Protected routes (auth required)
        protected := books.Group("")
        protected.Use(authMiddleware.AuthRequired())
        {
            protected.POST("", h.CreateBook)
            protected.PUT("/:id", h.UpdateBook)
            protected.DELETE("/:id", h.DeleteBook)
        }
    }
}
```

### Register in main router

```go
// internal/api/routers/router.go
// Add to SetupRouter function:

RegisterBookRoutes(v1, handler, authMiddleware)
```

---

## Step 7: Create Mock for Testing

Location: `tests/mocks/`

```go
// tests/mocks/book_repository_mock.go
package mocks

import (
    "context"

    "github.com/alireza-akbarzadeh/ginflow/internal/models"
    "github.com/alireza-akbarzadeh/ginflow/internal/query"
    "github.com/stretchr/testify/mock"
)

type MockBookRepository struct {
    mock.Mock
}

func (m *MockBookRepository) Insert(ctx context.Context, book *models.Book) (*models.Book, error) {
    args := m.Called(ctx, book)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookRepository) Get(ctx context.Context, id int) (*models.Book, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookRepository) GetByISBN(ctx context.Context, isbn string) (*models.Book, error) {
    args := m.Called(ctx, isbn)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookRepository) GetAll(ctx context.Context, req *query.QueryParams) ([]models.Book, *query.PaginatedList, error) {
    args := m.Called(ctx, mock.Anything)
    if args.Get(0) == nil {
        return nil, nil, args.Error(2)
    }
    return args.Get(0).([]models.Book), args.Get(1).(*query.PaginatedList), args.Error(2)
}

func (m *MockBookRepository) Update(ctx context.Context, book *models.Book) error {
    args := m.Called(ctx, book)
    return args.Error(0)
}

func (m *MockBookRepository) Delete(ctx context.Context, id int) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockBookRepository) GetByUser(ctx context.Context, userID int) ([]models.Book, error) {
    args := m.Called(ctx, userID)
    return args.Get(0).([]models.Book), args.Error(1)
}
```

---

## Step 8: Write Tests

Location: `tests/`

```go
// tests/book_test.go
package tests

import (
    "net/http"
    "testing"
    "time"

    "github.com/alireza-akbarzadeh/ginflow/internal/models"
    "github.com/alireza-akbarzadeh/ginflow/internal/query"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestBookManagement(t *testing.T) {
    suite := SetupMockTestSuite(t)

    t.Run("create book", func(t *testing.T) {
        // Setup mock expectation
        suite.MockBooks.On("GetByISBN", mock.Anything, "1234567890123").Return(nil, nil)
        suite.MockBooks.On("Insert", mock.Anything, mock.MatchedBy(func(b *models.Book) bool {
            return b.Title == "Test Book" && b.UserID == 1
        })).Return(&models.Book{
            ID:        1,
            UserID:    1,
            Title:     "Test Book",
            Author:    "Test Author",
            ISBN:      "1234567890123",
            CreatedAt: time.Now(),
        }, nil)

        // Make request
        body := `{"title":"Test Book","author":"Test Author","isbn":"1234567890123"}`
        w := suite.AuthenticatedRequest("POST", "/api/v1/books", body)

        // Assert
        assert.Equal(t, http.StatusCreated, w.Code)
    })

    t.Run("get all books", func(t *testing.T) {
        books := []models.Book{
            {ID: 1, Title: "Book 1", Author: "Author 1"},
            {ID: 2, Title: "Book 2", Author: "Author 2"},
        }

        result := &query.PaginatedList{
            Success: true,
            Data:    books,
            Pagination: &query.PageInfo{
                Page:        1,
                PageSize:    20,
                TotalItems:  2,
                TotalPages:  1,
                HasNextPage: false,
                HasPrevPage: false,
                Count:       2,
            },
        }

        suite.MockBooks.On("GetAll", mock.Anything, mock.Anything).Return(books, result, nil)

        w := suite.MakeRequest("GET", "/api/v1/books", "")

        assert.Equal(t, http.StatusOK, w.Code)
    })

    t.Run("get single book", func(t *testing.T) {
        book := &models.Book{ID: 1, Title: "Test Book", UserID: 1}
        suite.MockBooks.On("Get", mock.Anything, 1).Return(book, nil)

        w := suite.MakeRequest("GET", "/api/v1/books/1", "")

        assert.Equal(t, http.StatusOK, w.Code)
    })

    t.Run("book not found", func(t *testing.T) {
        suite.MockBooks.On("Get", mock.Anything, 999).Return(nil, nil)

        w := suite.MakeRequest("GET", "/api/v1/books/999", "")

        assert.Equal(t, http.StatusNotFound, w.Code)
    })

    t.Run("update own book", func(t *testing.T) {
        book := &models.Book{ID: 1, Title: "Old Title", UserID: 1}
        suite.MockBooks.On("Get", mock.Anything, 1).Return(book, nil)
        suite.MockBooks.On("Update", mock.Anything, mock.Anything).Return(nil)

        body := `{"title":"New Title","author":"New Author"}`
        w := suite.AuthenticatedRequest("PUT", "/api/v1/books/1", body)

        assert.Equal(t, http.StatusOK, w.Code)
    })

    t.Run("cannot update others book", func(t *testing.T) {
        book := &models.Book{ID: 1, Title: "Book", UserID: 999} // Different user
        suite.MockBooks.On("Get", mock.Anything, 1).Return(book, nil)

        body := `{"title":"Hacked Title"}`
        w := suite.AuthenticatedRequest("PUT", "/api/v1/books/1", body)

        assert.Equal(t, http.StatusForbidden, w.Code)
    })

    t.Run("delete own book", func(t *testing.T) {
        book := &models.Book{ID: 1, UserID: 1}
        suite.MockBooks.On("Get", mock.Anything, 1).Return(book, nil)
        suite.MockBooks.On("Delete", mock.Anything, 1).Return(nil)

        w := suite.AuthenticatedRequest("DELETE", "/api/v1/books/1", "")

        assert.Equal(t, http.StatusNoContent, w.Code)
    })
}
```

---

## Error Handling

### Error Types

| Error Type         | HTTP Status | When to Use            |
| ------------------ | ----------- | ---------------------- |
| `ErrNotFound`      | 404         | Resource doesn't exist |
| `ErrAlreadyExists` | 409         | Duplicate resource     |
| `ErrForbidden`     | 403         | No permission          |
| `ErrUnauthorized`  | 401         | Not authenticated      |
| `ErrValidation`    | 400         | Invalid input          |
| `ErrInternal`      | 500         | Server error           |

### Using Error Helpers

```go
import appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"

// Simple error
helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrNotFound, "book not found"), "")

// Formatted error
helpers.RespondWithAppError(c, appErrors.Newf(appErrors.ErrNotFound, "book with ID %d not found", id), "")

// Handle repository errors
if helpers.HandleError(c, err, "Failed to create book") {
    return // Error response already sent
}
```

### Error Response Format

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "book with ID 123 not found"
  }
}
```

---

## Pagination Guide

### Query Parameters

| Parameter   | Type   | Description                            | Example                             |
| ----------- | ------ | -------------------------------------- | ----------------------------------- |
| `page`      | int    | Page number (default: 1)               | `?page=2`                           |
| `page_size` | int    | Items per page (default: 20, max: 100) | `?page_size=50`                     |
| `type`      | string | `offset` or `cursor`                   | `?type=cursor`                      |
| `cursor`    | string | Cursor for cursor pagination           | `?cursor=abc123`                    |
| `sort`      | string | Sort fields                            | `?sort=-created_at,name:asc`        |
| `search`    | string | Search term                            | `?search=golang`                    |
| `field[op]` | any    | Filter by field                        | `?price[gte]=100&status[eq]=active` |

### Sort Syntax

```
# Descending (prefix with -)
?sort=-created_at

# Ascending (no prefix or :asc)
?sort=name
?sort=name:asc

# Multiple fields
?sort=-created_at,name:asc,price:desc
```

### Filter Operators

| Operator  | Description                 | Example                      |
| --------- | --------------------------- | ---------------------------- |
| `eq`      | Equal                       | `?status[eq]=active`         |
| `neq`     | Not equal                   | `?status[neq]=deleted`       |
| `gt`      | Greater than                | `?price[gt]=100`             |
| `gte`     | Greater or equal            | `?price[gte]=100`            |
| `lt`      | Less than                   | `?price[lt]=1000`            |
| `lte`     | Less or equal               | `?price[lte]=1000`           |
| `like`    | Contains                    | `?name[like]=book`           |
| `ilike`   | Contains (case-insensitive) | `?name[ilike]=BOOK`          |
| `in`      | In list                     | `?status[in]=active,pending` |
| `null`    | Is null                     | `?deleted_at[null]`          |
| `notnull` | Is not null                 | `?image[notnull]`            |

### Response Format

```json
{
    "data": [...],
    "pagination": {
        "current_page": 1,
        "page_size": 20,
        "total_items": 100,
        "total_pages": 5,
        "has_next": true,
        "has_prev": false,
        "next_cursor": "eyJpZCI6MjB9",
        "prev_cursor": ""
    },
    "links": {
        "self": "/api/v1/books?page=1&page_size=20",
        "first": "/api/v1/books?page=1&page_size=20",
        "last": "/api/v1/books?page=5&page_size=20",
        "next": "/api/v1/books?page=2&page_size=20",
        "prev": ""
    }
}
```

---

## Validation Guide

### Model Validation Tags

```go
type Book struct {
    Title       string  `binding:"required,min=2,max=255"`    // Required, 2-255 chars
    Price       float64 `binding:"omitempty,gte=0"`           // Optional, >= 0
    Email       string  `binding:"required,email"`            // Required, valid email
    Status      string  `binding:"required,oneof=draft published"` // Must be one of
    ISBN        string  `binding:"omitempty,len=13"`          // Optional, exactly 13 chars
    URL         string  `binding:"omitempty,url"`             // Optional, valid URL
    Rating      int     `binding:"omitempty,min=1,max=5"`     // Optional, 1-5
}
```

### Custom Validation

```go
// In handler
if book.Price < 0 {
    helpers.RespondWithError(c, http.StatusBadRequest, "price cannot be negative")
    return
}
```

---

## Best Practices

### 1. Always Use Context

```go
// ✅ Good
r.DB.WithContext(ctx).Find(&items)

// ❌ Bad
r.DB.Find(&items)
```

### 2. Log Important Operations

```go
logging.Debug(ctx, "starting operation", "param", value)
logging.Info(ctx, "operation completed", "result_id", id)
logging.Error(ctx, "operation failed", "error", err)
```

### 3. Check for nil Before Use

```go
// ✅ Good
if book == nil {
    helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrNotFound, "not found"), "")
    return
}
// Now safe to use book

// ❌ Bad
c.JSON(http.StatusOK, book) // Might be nil!
```

### 4. Validate Ownership

```go
if resource.UserID != user.ID {
    helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrForbidden, "not allowed"), "")
    return
}
```

### 5. Use Transactions for Multiple Operations

```go
func (r *Repository) CreateWithRelations(ctx context.Context, item *Model) error {
    return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(item).Error; err != nil {
            return err
        }
        if err := tx.Create(&RelatedItem{ItemID: item.ID}).Error; err != nil {
            return err
        }
        return nil
    })
}
```

### 6. Whitelist Filter/Sort Fields

```go
builder := query.NewQueryBuilder(db).
    AllowFilters("name", "status", "price").  // Only these can be filtered
    AllowSorts("name", "created_at")          // Only these can be sorted
```

---

## Common Mistakes to Avoid

### ❌ Don't: Return error for not found

```go
// Bad
if result.Error == gorm.ErrRecordNotFound {
    return nil, result.Error  // Don't return error
}
```

### ✅ Do: Return nil, nil for not found

```go
// Good
if result.Error == gorm.ErrRecordNotFound {
    return nil, nil  // Not found is not an error
}
```

### ❌ Don't: Skip error handling

```go
// Bad
book, _ := h.Repos.Books.Get(ctx, id)
c.JSON(200, book)
```

### ✅ Do: Always handle errors

```go
// Good
book, err := h.Repos.Books.Get(ctx, id)
if helpers.HandleError(c, err, "Failed to get book") {
    return
}
if book == nil {
    helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrNotFound, "not found"), "")
    return
}
c.JSON(200, book)
```

### ❌ Don't: Allow arbitrary filter/sort fields

```go
// Bad - SQL injection risk
query.Where(fmt.Sprintf("%s = ?", userInput), value)
```

### ✅ Do: Whitelist allowed fields

```go
// Good
builder.AllowFilters("name", "status")  // Only whitelisted fields
```

### ❌ Don't: Create duplicate code

```go
// Bad - repeated in every handler
id, err := strconv.Atoi(c.Param("id"))
if err != nil {
    c.JSON(400, gin.H{"error": "invalid id"})
    return
}
```

### ✅ Do: Use helpers

```go
// Good
id, err := helpers.ParseIDParam(c, "id")
if err != nil {
    return  // Helper already sent response
}
```

### ❌ Don't: Forget authentication check

```go
// Bad
func (h *Handler) DeleteBook(c *gin.Context) {
    id, _ := helpers.ParseIDParam(c, "id")
    h.Repos.Books.Delete(ctx, id)  // Anyone can delete!
}
```

### ✅ Do: Always verify ownership

```go
// Good
func (h *Handler) DeleteBook(c *gin.Context) {
    user, ok := helpers.GetAuthenticatedUser(c)
    if !ok {
        return
    }
    book, _ := h.Repos.Books.Get(ctx, id)
    if book.UserID != user.ID {
        helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrForbidden, "not allowed"), "")
        return
    }
    h.Repos.Books.Delete(ctx, id)
}
```

---

## File Structure Summary

```
internal/
├── models/
│   └── book.go              # 1. Model definition
├── repository/
│   ├── interfaces/
│   │   └── book.go          # 2. Repository interface
│   ├── book.go              # 3. Repository implementation
│   └── repository.go        # 4. Register repository
├── api/
│   ├── handlers/
│   │   └── book.go          # 5. HTTP handlers
│   └── routers/
│       ├── book.go          # 6. Route definitions
│       └── router.go        # 6. Register routes
tests/
├── mocks/
│   └── book_repository_mock.go  # 7. Test mock
└── book_test.go             # 8. Unit tests
```

---

## Quick Reference Checklist

When creating a new API resource:

- [ ] Create model in `internal/models/`
- [ ] Add GORM tags for database
- [ ] Add binding tags for validation
- [ ] Create interface in `internal/repository/interfaces/`
- [ ] Implement repository in `internal/repository/`
- [ ] Include `GetAll` method with `QueryParams` support
- [ ] Register in `internal/repository/repository.go`
- [ ] Create handlers in `internal/api/handlers/`
- [ ] Add Swagger annotations
- [ ] Use error helpers consistently
- [ ] Add logging for operations
- [ ] Create router in `internal/api/routers/`
- [ ] Register routes in main router
- [ ] Create mock in `tests/mocks/`
- [ ] Write unit tests in `tests/`
- [ ] Run tests: `go test ./tests/... -v`
- [ ] Update Swagger: `swag init -g cmd/server/main.go`

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/alireza-akbarzadeh/ginflow/internal/api/routers"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/repository"
	"github.com/alireza-akbarzadeh/ginflow/tests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestSuite holds the test application and database
type TestSuite struct {
	DB        *gorm.DB
	Router    *gin.Engine
	Handler   *handlers.Handler
	JWTSecret string
	Mocks     *repository.Models
}

// SetupMockTestSuite initializes the test suite with mocks
func SetupMockTestSuite(t *testing.T) *TestSuite {
	gin.SetMode(gin.TestMode)

	// Initialize logging for tests
	logging.InitLogger()

	// Create mocks
	mockRepos := &repository.Models{
		Users:      &mocks.UserRepositoryMock{},
		Events:     &mocks.EventRepositoryMock{},
		Attendees:  &mocks.AttendeeRepositoryMock{},
		Categories: &mocks.CategoryRepositoryMock{},
		Comments:   &mocks.CommentRepositoryMock{},
		Profiles:   &mocks.ProfileRepositoryMock{},
		Products:   &mocks.ProductRepositoryMock{},
		Baskets:    &mocks.BasketRepositoryMock{},
	}

	// JWT secret for testing
	jwtSecret := "test-jwt-secret-key"

	// Create handler
	handler := handlers.NewHandler(mockRepos, jwtSecret)

	// Create router
	router := routers.SetupRouter(handler, jwtSecret, mockRepos.Users)

	return &TestSuite{
		Router:    router,
		Handler:   handler,
		JWTSecret: jwtSecret,
		Mocks:     mockRepos,
	}
}

// SetupTestSuite initializes the test suite with database and router
func SetupTestSuite(t *testing.T) *TestSuite {
	gin.SetMode(gin.TestMode)

	// Load environment variables from .env file
	if err := godotenv.Load("../.env"); err != nil {
		t.Logf("Warning: Could not load .env file: %v", err)
	}

	// Get database URL from environment variable (prefer TEST_DATABASE_URL for tests)
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		// If neither DATABASE_URL nor TEST_DATABASE_URL is set, skip tests
		t.Skip("DATABASE_URL or TEST_DATABASE_URL environment variable is required for testing, skipping...")
		return nil
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// If database connection fails, skip tests
		t.Skip("Failed to connect to database for testing, skipping...")
		return nil
	}

	// Auto-migrate all models
	err = db.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Attendee{},
		&models.Category{},
		&models.Comment{},
	)
	require.NoError(t, err)

	// Create repositories
	models := repository.NewModels(db)

	// JWT secret for testing
	jwtSecret := "test-jwt-secret-key"

	// Create handler
	handler := handlers.NewHandler(models, jwtSecret)

	// Create router
	router := routers.SetupRouter(handler, jwtSecret, models.Users)

	return &TestSuite{
		DB:        db,
		Router:    router,
		Handler:   handler,
		JWTSecret: jwtSecret,
	}
}

// TeardownTestSuite cleans up the test suite
func (ts *TestSuite) TeardownTestSuite(t *testing.T) {
	sqlDB, err := ts.DB.DB()
	require.NoError(t, err)
	sqlDB.Close()
}

// Helper function to create authenticated request
func (ts *TestSuite) createAuthenticatedRequest(method, path, token string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	ts.Router.ServeHTTP(w, req)
	return w
}

// Helper function to create unauthenticated request
func (ts *TestSuite) createRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	return ts.createAuthenticatedRequest(method, path, "", body)
}

// Helper function to register and login a test user
func (ts *TestSuite) createTestUser(t *testing.T, email, password, name string) (string, *models.User) {
	// Try to register user (may fail if user already exists)
	registerReq := handlers.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
	}
	w := ts.createRequest("POST", "/api/v1/auth/register", registerReq)
	// Accept both 201 (created) and 409 (already exists)
	if w.Code != http.StatusCreated && w.Code != http.StatusConflict {
		t.Fatalf("Expected status 201 or 409, got %d", w.Code)
	}

	// Login to get token and user data
	loginReq := handlers.LoginRequest{
		Email:    email,
		Password: password,
	}
	w = ts.createRequest("POST", "/api/v1/auth/login", loginReq)
	assert.Equal(t, http.StatusOK, w.Code)

	var loginResp handlers.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResp)
	require.NoError(t, err)

	// Get user data from login response
	var user models.User
	user.ID = loginResp.User.ID
	user.Email = loginResp.User.Email
	user.Name = loginResp.User.Name

	return loginResp.Token, &user
}

// Helper function to create a test event
func (ts *TestSuite) createTestEvent(t *testing.T, token string, event models.Event) *models.Event {
	w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, event)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdEvent models.Event
	err := json.Unmarshal(w.Body.Bytes(), &createdEvent)
	require.NoError(t, err)

	return &createdEvent
}

// GenerateToken generates a JWT token for testing
func (ts *TestSuite) GenerateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(ts.JWTSecret))
}

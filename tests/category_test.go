package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCategoryManagement tests category CRUD operations
func TestCategoryManagement(t *testing.T) {
	ts := SetupMockTestSuite(t)

	// Create test user for authentication
	userID := 1
	token, err := ts.GenerateToken(userID)
	assert.NoError(t, err)

	// Mock user retrieval for authentication
	mockUserRepo := ts.Mocks.Users.(*mocks.UserRepositoryMock)
	mockUserRepo.On("Get", mock.Anything, userID).Return(&models.User{
		ID:    userID,
		Email: "categoryuser@example.com",
		Name:  "Category User",
	}, nil)

	mockCategoryRepo := ts.Mocks.Categories.(*mocks.CategoryRepositoryMock)

	t.Run("create category", func(t *testing.T) {
		category := models.Category{
			Name:        "Technology",
			Description: "Technology-related events",
		}

		// Expect Insert call
		mockCategoryRepo.On("Insert", mock.Anything, mock.MatchedBy(func(c *models.Category) bool {
			return c.Name == category.Name && c.Description == category.Description
		})).Return(&models.Category{ID: 1, Name: category.Name, Description: category.Description}, nil).Once()

		w := ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, category)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdCategory models.Category
		err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
		assert.NoError(t, err)
		assert.Equal(t, category.Name, createdCategory.Name)
		assert.Equal(t, category.Description, createdCategory.Description)
	})

	t.Run("get all categories", func(t *testing.T) {
		categories := []*models.Category{
			{Name: "Sports", Description: "Sports events"},
			{Name: "Music", Description: "Music concerts and events"},
		}

		// Expect GetAll calls
		mockCategoryRepo.On("GetAll", mock.Anything).Return(categories, nil).Once()

		// Get all categories
		w := ts.createRequest("GET", "/api/v1/categories", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var respCategories []models.Category
		err := json.Unmarshal(w.Body.Bytes(), &respCategories)
		assert.NoError(t, err)
		assert.Equal(t, len(categories), len(respCategories))
		assert.Equal(t, categories[0].Name, respCategories[0].Name)
		assert.Equal(t, categories[1].Name, respCategories[1].Name)
	})

	t.Run("category validation", func(t *testing.T) {
		// Test empty name
		invalidCategory := models.Category{Name: "", Description: "Valid description"}
		w := ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, invalidCategory)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test name too short
		shortNameCategory := models.Category{Name: "AB", Description: "Valid description"}
		w = ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, shortNameCategory)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("get categories without authentication", func(t *testing.T) {
		mockCategoryRepo.On("GetAll", mock.Anything).Return([]*models.Category{}, nil).Once()

		// Getting categories should work without authentication
		w := ts.createRequest("GET", "/api/v1/categories", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var categories []models.Category
		err := json.Unmarshal(w.Body.Bytes(), &categories)
		assert.NoError(t, err)
		assert.IsType(t, []models.Category{}, categories)
	})

	t.Run("create category requires authentication", func(t *testing.T) {
		category := models.Category{
			Name:        "Unauthenticated Category",
			Description: "This should fail",
		}

		w := ts.createRequest("POST", "/api/v1/categories", category)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// TestCategoryEdgeCases tests edge cases for categories
func TestCategoryEdgeCases(t *testing.T) {
	ts := SetupMockTestSuite(t)

	userID := 2
	token, err := ts.GenerateToken(userID)
	assert.NoError(t, err)

	mockUserRepo := ts.Mocks.Users.(*mocks.UserRepositoryMock)
	mockUserRepo.On("Get", mock.Anything, userID).Return(&models.User{
		ID:    userID,
		Email: "edgecategoryuser@example.com",
		Name:  "Edge Category User",
	}, nil)

	mockCategoryRepo := ts.Mocks.Categories.(*mocks.CategoryRepositoryMock)

	t.Run("category with empty description", func(t *testing.T) {
		category := models.Category{
			Name:        "Empty Description Category",
			Description: "",
		}

		mockCategoryRepo.On("Insert", mock.Anything, mock.MatchedBy(func(c *models.Category) bool {
			return c.Name == category.Name && c.Description == ""
		})).Return(&models.Category{ID: 2, Name: category.Name, Description: ""}, nil).Once()

		w := ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, category)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdCategory models.Category
		err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
		assert.NoError(t, err)
		assert.Equal(t, category.Name, createdCategory.Name)
		assert.Equal(t, "", createdCategory.Description)
	})

	t.Run("category name with special characters", func(t *testing.T) {
		category := models.Category{
			Name:        "Tech & Innovation",
			Description: "Category with special characters",
		}

		mockCategoryRepo.On("Insert", mock.Anything, mock.MatchedBy(func(c *models.Category) bool {
			return c.Name == category.Name
		})).Return(&models.Category{ID: 3, Name: category.Name}, nil).Once()

		w := ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, category)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdCategory models.Category
		err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
		assert.NoError(t, err)
		assert.Equal(t, "Tech & Innovation", createdCategory.Name)
	})

	t.Run("very long category name", func(t *testing.T) {
		longName := "This is a very long category name that exceeds normal length but should still be accepted since we don't have explicit length limits in validation"
		category := models.Category{
			Name:        longName,
			Description: "Category with very long name",
		}

		mockCategoryRepo.On("Insert", mock.Anything, mock.MatchedBy(func(c *models.Category) bool {
			return c.Name == longName
		})).Return(&models.Category{ID: 4, Name: longName}, nil).Once()

		w := ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, category)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdCategory models.Category
		err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
		assert.NoError(t, err)
		assert.Equal(t, longName, createdCategory.Name)
	})
}

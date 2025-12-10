package handlers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/helpers"
	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/utils"
	"github.com/gin-gonic/gin"
)

// CreateCategory handles category creation
// @Summary      Create a new category
// @Description  Create a new event category (requires authentication)
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        category  body      models.Category  true  "Category object"
// @Success      201       {object}  models.Category
// @Failure      400       {object}  helpers.ErrorResponse
// @Failure      401       {object}  helpers.ErrorResponse
// @Failure      500       {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/categories [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	ctx := c.Request.Context()

	var category models.Category
	if !helpers.BindJSON(c, &category) {
		return
	}

	// Generate slug if not provided
	if category.Slug == "" {
		category.Slug = utils.GenerateSlug(category.Name)
	}

	createdCategory, err := h.Repos.Categories.Insert(ctx, &category)
	if err != nil {
		logging.Error(ctx, "Failed to create category", err, "name", category.Name)
		helpers.HandleError(c, err, "Failed to create category")
		return
	}

	c.JSON(http.StatusCreated, createdCategory)
}

// GetAllCategories retrieves all categories
// @Summary      Get all categories
// @Description  Get a list of all event categories
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Category
// @Failure      500  {object}  helpers.ErrorResponse
// @Router       /api/v1/categories [get]
func (h *Handler) GetAllCategories(c *gin.Context) {
	ctx := c.Request.Context()

	categories, err := h.Repos.Categories.GetAll(ctx)
	if err != nil {
		logging.Error(ctx, "Failed to retrieve categories", err)
		helpers.HandleError(c, err, "Failed to retrieve categories")
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryBySlug retrieves a category by slug
// @Summary      Get category by slug
// @Description  Get a category by its slug
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        slug   path      string  true  "Category Slug"
// @Success      200    {object}  models.Category
// @Failure      404    {object}  helpers.ErrorResponse
// @Failure      500    {object}  helpers.ErrorResponse
// @Router       /api/v1/categories/{slug} [get]
func (h *Handler) GetCategoryBySlug(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	category, err := h.Repos.Categories.GetBySlug(ctx, slug)
	if err != nil {
		logging.Error(ctx, "Failed to retrieve category", err, "slug", slug)
		helpers.HandleError(c, err, "Failed to retrieve category")
		return
	}
	if category == nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrNotFound, "Category not found"), "Category not found")
		return
	}

	c.JSON(http.StatusOK, category)
}

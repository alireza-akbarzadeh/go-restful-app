package handlers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/ginflow/pkg/api/helpers"
	"github.com/alireza-akbarzadeh/ginflow/pkg/repository"
	"github.com/gin-gonic/gin"
)

// CreateCategory handles category creation
// @Summary      Create a new category
// @Description  Create a new event category (requires authentication)
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        category  body      repository.Category  true  "Category object"
// @Success      201       {object}  repository.Category
// @Failure      400       {object}  helpers.ErrorResponse
// @Failure      401       {object}  helpers.ErrorResponse
// @Failure      500       {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/categories [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	var category repository.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	createdCategory, err := h.Repos.Categories.Insert(&category)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to create category")
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
// @Success      200  {array}   repository.Category
// @Failure      500  {object}  helpers.ErrorResponse
// @Router       /api/v1/categories [get]
func (h *Handler) GetAllCategories(c *gin.Context) {
	categories, err := h.Repos.Categories.GetAll()
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve categories")
		return
	}

	c.JSON(http.StatusOK, categories)
}

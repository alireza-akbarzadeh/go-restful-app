package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupCategoryRoutes configures public category routes
func SetupCategoryRoutes(router *gin.RouterGroup, h *handlers.Handler) {
	categories := router.Group("/categories")
	{
		categories.GET("", h.GetAllCategories)
		categories.GET("/:slug", h.GetCategoryBySlug)
	}
}

// SetupProtectedCategoryRoutes configures protected category routes
func SetupProtectedCategoryRoutes(router *gin.RouterGroup, h *handlers.Handler) {
	router.POST("/categories", h.CreateCategory)
}

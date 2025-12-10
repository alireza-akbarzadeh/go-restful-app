package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupProductRoutes configures product-related routes
func SetupProductRoutes(router *gin.RouterGroup, h *handlers.Handler) {
	products := router.Group("/products")
	{
		products.GET("", h.GetAllProducts)
		products.GET("/:id", h.GetProduct)
		products.GET("/slug/:slug", h.GetProductBySlug)
		products.GET("/category/:id", h.GetProductsByCategory)
	}
}

// SetupProtectedProductRoutes configures protected product-related routes
func SetupProtectedProductRoutes(router *gin.RouterGroup, h *handlers.Handler) {
	products := router.Group("/products")
	{
		products.POST("", h.CreateProduct)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
	}
}

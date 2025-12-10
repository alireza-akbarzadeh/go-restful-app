package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupProtectedBasketRoutes registers protected basket routes
func SetupProtectedBasketRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	basket := rg.Group("/basket")
	{
		basket.GET("", h.GetBasket)
		basket.DELETE("", h.ClearBasket)
		basket.POST("/items", h.AddItemToBasket)
		basket.DELETE("/items/:id", h.RemoveItemFromBasket)
	}
}

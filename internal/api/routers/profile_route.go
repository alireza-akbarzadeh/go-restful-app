package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupProtectedProfileRoutes configures protected profile routes
func SetupProtectedProfileRoutes(router *gin.RouterGroup, h *handlers.Handler) {
	router.GET("/profile", h.GetProfile)
	router.POST("/profile", h.CreateProfile)
	router.PUT("/profile", h.UpdateProfile)
	router.DELETE("/profile", h.DeleteProfile)
}

package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes configures authentication-related routes
func SetupAuthRoutes(router *gin.RouterGroup, h *handlers.Handler) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
	}
}

// SetupProtectedAuthRoutes configures protected authentication-related routes
func SetupProtectedAuthRoutes(router *gin.RouterGroup, h *handlers.Handler) {
	auth := router.Group("/auth")
	{
		auth.PUT("/password", h.UpdatePassword)
	}
}

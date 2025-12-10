package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupProtectedUserRoutes configures protected user routes
func SetupProtectedUserRoutes(router *gin.RouterGroup, h *handlers.Handler) {
	router.GET("/users", h.GetAllUsers)
	router.PUT("/users/:id", h.UpdateUser)
	router.DELETE("/users/:id", h.DeleteUser)
}

package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupEventRoutes configures public event routes
func SetupEventRoutes(router *gin.RouterGroup, h *handlers.Handler) {
	events := router.Group("/events")
	{
		events.GET("", h.GetAllEvents)
		events.GET("/:id", h.GetEvent)
		events.GET("/:id/attendees", h.GetAttendees)
		events.GET("/:id/comments", h.GetEventComments)
	}
}

// SetupProtectedEventRoutes configures protected event routes
func SetupProtectedEventRoutes(router *gin.RouterGroup, h *handlers.Handler) {
	// Event management
	router.POST("/events", h.CreateEvent)
	router.PUT("/events/:id", h.UpdateEvent)
	router.DELETE("/events/:id", h.DeleteEvent)

	// Comment management
	router.POST("/events/:id/comments", h.CreateComment)
	router.DELETE("/events/:id/comments/:commentId", h.DeleteComment)

	// Attendee management
	router.POST("/events/:id/attendees/:userId", h.AddAttendee)
	router.DELETE("/events/:id/attendees/:userId", h.RemoveAttendee)
}

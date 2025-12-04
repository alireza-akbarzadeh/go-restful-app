package routers

import (
	"github.com/alireza-akbarzadeh/restful-app/pkg/api/handlers"
	"github.com/alireza-akbarzadeh/restful-app/pkg/api/middleware"
	"github.com/alireza-akbarzadeh/restful-app/pkg/repository"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures and returns the main router
func SetupRouter(handler *handlers.Handler, jwtSecret string, userRepo *repository.UserRepository) *gin.Engine {
	router := gin.Default()

	// Setup Swagger documentation
	router.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(302, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}

		// Public event routes
		events := v1.Group("/events")
		{
			events.GET("", handler.GetAllEvents)
			events.GET("/:id", handler.GetEvent)
			events.GET("/:id/attendees", handler.GetAttendees)
		}

		// Public attendee routes
		attendees := v1.Group("/attendees")
		{
			attendees.GET("/:id/events", handler.GetEventsByAttendee)
		}

		// Public category routes
		categories := v1.Group("/categories")
		{
			categories.GET("", handler.GetAllCategories)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret, userRepo))
		{
			// Event management
			protected.POST("/events", handler.CreateEvent)
			protected.PUT("/events/:id", handler.UpdateEvent)
			protected.DELETE("/events/:id", handler.DeleteEvent)

			// Attendee management
			protected.POST("/events/:id/attendees/:userId", handler.AddAttendee)
			protected.DELETE("/events/:id/attendees/:userId", handler.RemoveAttendee)

			// Category management
			protected.POST("/categories", handler.CreateCategory)
		}
	}

	return router
}

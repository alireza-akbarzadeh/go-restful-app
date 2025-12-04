package routers

import (
	"html/template"

	"github.com/alireza-akbarzadeh/restful-app/pkg/api/handlers"
	"github.com/alireza-akbarzadeh/restful-app/pkg/api/middleware"
	"github.com/alireza-akbarzadeh/restful-app/pkg/repository"
	"github.com/alireza-akbarzadeh/restful-app/pkg/web"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures and returns the main router
func SetupRouter(handler *handlers.Handler, jwtSecret string, userRepo *repository.UserRepository) *gin.Engine {
	router := gin.Default()
	router.SetHTMLTemplate(template.Must(template.ParseFS(web.Templates, "*.html")))

	// Root landing page
	router.GET("/", handler.ShowLandingPage)

	// Setup Swagger documentation
	router.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(302, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	// Health check endpoint
	router.GET("/health", handler.ShowHealthPage)

	// Dashboard endpoint
	router.GET("/dashboard", handler.ShowDashboardPage)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
			auth.POST("/logout", handler.Logout)
		}

		// Public event routes
		events := v1.Group("/events")
		{
			events.GET("", handler.GetAllEvents)
			events.GET("/:id", handler.GetEvent)
			events.GET("/:id/attendees", handler.GetAttendees)
			events.GET("/:id/comments", handler.GetEventComments)
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

			// Comment management
			protected.POST("/events/:id/comments", handler.CreateComment)
			protected.DELETE("/events/:id/comments/:commentId", handler.DeleteComment)

			// Attendee management
			protected.POST("/events/:id/attendees/:userId", handler.AddAttendee)
			protected.DELETE("/events/:id/attendees/:userId", handler.RemoveAttendee)

			// Category management
			protected.POST("/categories", handler.CreateCategory)

			// User management
			protected.PUT("/auth/password", handler.UpdatePassword)
			protected.GET("/users", handler.GetAllUsers)
			protected.PUT("/users/:id", handler.UpdateUser)
			protected.DELETE("/users/:id", handler.DeleteUser)
		}
	}

	return router
}

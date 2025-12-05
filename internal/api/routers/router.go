package routers

import (
	"html/template"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/alireza-akbarzadeh/ginflow/internal/api/middleware"
	"github.com/alireza-akbarzadeh/ginflow/internal/constants"
	"github.com/alireza-akbarzadeh/ginflow/internal/repository/interfaces"
	"github.com/alireza-akbarzadeh/ginflow/internal/web"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

// SetupRouter configures and returns the main router
func SetupRouter(handler *handlers.Handler, jwtSecret string, userRepo interfaces.UserRepositoryInterface) *gin.Engine {
	// Use gin.New() instead of Default() to use custom logger
	router := gin.New()

	// Add Recovery middleware to recover from panics
	router.Use(gin.Recovery())

	// Add Structured Logger middleware
	router.Use(middleware.Logger())

	// Apply rate limiting middleware globally
	// 20 requests per second with a burst of 50
	router.Use(middleware.RateLimitMiddleware(rate.Limit(constants.DEFAULT_RATE_LIMIT), constants.DEFAULT_RATE_BURST))

	// Apply CORS middleware to allow requests from frontend
	router.Use(middleware.CORS([]string{"*"}))

	// Apply security headers middleware
	router.Use(middleware.SecurityHeaders())

	router.SetHTMLTemplate(template.Must(template.ParseFS(web.Templates, "components/*.html", "pages/*.html")))

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
		// Auth Routes
		SetupAuthRoutes(v1, handler)

		// Event Routes
		SetupEventRoutes(v1, handler)

		// Attendee Routes
		SetupAttendeeRoutes(v1, handler)

		// Category Routes
		SetupCategoryRoutes(v1, handler)

		// Product Routes
		SetupProductRoutes(v1, handler)

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret, userRepo))
		{
			SetupProtectedAuthRoutes(protected, handler)
			SetupProtectedEventRoutes(protected, handler)
			SetupProtectedCategoryRoutes(protected, handler)
			SetupProtectedUserRoutes(protected, handler)
			SetupProtectedProfileRoutes(protected, handler)
			SetupProtectedProductRoutes(protected, handler)
			SetupProtectedBasketRoutes(protected, handler)

		}
	}

	return router
}

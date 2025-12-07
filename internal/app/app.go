package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/alireza-akbarzadeh/ginflow/internal/api/routers"
	"github.com/alireza-akbarzadeh/ginflow/internal/console"
	"github.com/alireza-akbarzadeh/ginflow/internal/database"
	"github.com/alireza-akbarzadeh/ginflow/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// App is the main application container that orchestrates all components
type App struct {
	config  *Config
	console *console.Console
	db      *gorm.DB
	sqlDB   *sql.DB
	repos   *repository.Models
	handler *handlers.Handler
	router  *gin.Engine
	server  *http.Server
}

// New creates a new App instance with the given options
func New(opts ...Option) (*App, error) {
	app := &App{
		config:  DefaultConfig(),
		console: console.New(),
	}

	// Apply functional options
	for _, opt := range opts {
		opt(app)
	}

	// Setup structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize all components
	if err := app.initialize(); err != nil {
		return nil, err
	}

	return app, nil
}

// initialize sets up all application components
func (a *App) initialize() error {
	a.console.Line()
	a.console.Info("🔧", "Initializing GinFlow Application...")
	a.console.Line()

	// 1. Validate Configuration
	a.console.Info("⚙️ ", "Loading configuration...")
	if a.config.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	a.console.Success("✅", "Configuration loaded")

	// 2. Initialize Database
	a.console.Info("🗄️ ", "Connecting to database...")
	db, err := database.Connect(a.config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	a.db = db

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}
	a.sqlDB = sqlDB
	a.console.Success("✅", "Database connected")

	// 3. Initialize Repositories
	a.console.Info("🔗", "Initializing dependencies...")
	a.repos = repository.NewModels(a.db)

	// 4. Initialize Handlers
	a.handler = handlers.NewHandler(a.repos, a.config.JWTSecret)

	// 5. Initialize Router
	a.router = routers.SetupRouter(a.handler, a.config.JWTSecret, a.repos.Users)
	a.console.Success("✅", "Dependencies initialized")

	// 6. Configure HTTP Server
	a.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", a.config.Port),
		Handler:      a.router,
		IdleTimeout:  a.config.IdleTimeout,
		ReadTimeout:  a.config.ReadTimeout,
		WriteTimeout: a.config.WriteTimeout,
	}

	return nil
}

// DB returns the GORM database instance
func (a *App) DB() *gorm.DB {
	return a.db
}

// Router returns the Gin router
func (a *App) Router() *gin.Engine {
	return a.router
}

// Config returns the application configuration
func (a *App) Config() *Config {
	return a.config
}

// Repos returns the repository models
func (a *App) Repos() *repository.Models {
	return a.repos
}

// Handler returns the HTTP handler
func (a *App) Handler() *handlers.Handler {
	return a.handler
}

// Server returns the HTTP server
func (a *App) Server() *http.Server {
	return a.server
}

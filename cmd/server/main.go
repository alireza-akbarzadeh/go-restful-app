package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/alireza-akbarzadeh/ginflow/docs"
	"github.com/alireza-akbarzadeh/ginflow/pkg/api/handlers"
	"github.com/alireza-akbarzadeh/ginflow/pkg/api/routers"
	"github.com/alireza-akbarzadeh/ginflow/pkg/config"
	"github.com/alireza-akbarzadeh/ginflow/pkg/database"
	"github.com/alireza-akbarzadeh/ginflow/pkg/repository"
	_ "github.com/joho/godotenv/autoload"
)

// @title Go Gin REST API
// @version 1.0
// @description This is a REST API server for managing events and attendees.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 0. Setup Structured Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 1. Load Configuration
	port := config.GetEnvInt("PORT", 8080)
	jwtSecret := config.GetEnvString("JWT_SECRET", "some-secret-123456")
	dbUrl := config.GetEnvString("DATABASE_URL", "")

	// 2. Initialize Database
	db, err := database.Connect(dbUrl)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}

	// Get underlying SQL DB for graceful shutdown
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get SQL DB", "error", err)
		os.Exit(1)
	}
	defer func() {
		slog.Info("Closing database connection...")
		if err := sqlDB.Close(); err != nil {
			slog.Error("Error closing database connection", "error", err)
		}
	}()

	// 3. Run Migrations
	if err := database.Migrate(db); err != nil {
		slog.Error("Migration failed", "error", err)
		os.Exit(1)
	}

	// 4. Initialize Dependencies
	repos := repository.NewModels(db)
	handler := handlers.NewHandler(repos, jwtSecret)
	router := routers.SetupRouter(handler, jwtSecret, repos.Users)

	// 5. Configure Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 6. Start Server in a Goroutine
	go func() {
		slog.Info("🚀 Server starting", "port", port)
		slog.Info("📚 Swagger UI", "url", fmt.Sprintf("http://localhost:%d/swagger/index.html", port))
		slog.Info("🏥 Health Check", "url", fmt.Sprintf("http://localhost:%d/health", port))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// 7. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")
}

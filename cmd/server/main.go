package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/alireza-akbarzadeh/restful-app/docs"
	"github.com/alireza-akbarzadeh/restful-app/pkg/api/handlers"
	"github.com/alireza-akbarzadeh/restful-app/pkg/api/routers"
	"github.com/alireza-akbarzadeh/restful-app/pkg/config"
	"github.com/alireza-akbarzadeh/restful-app/pkg/repository"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	// Load configuration
	port := config.GetEnvInt("PORT", 8080)
	jwtSecret := config.GetEnvString("JWT_SECRET", "some-secret-123456")
	dbUrl := config.GetEnvString("DATABASE_URL", "")

	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	log.Println("Using PostgreSQL database")
	dialector := postgres.Open(dbUrl)

	// Initialize database with GORM
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Auto Migrate the schema
	err = db.AutoMigrate(&repository.User{}, &repository.Event{}, &repository.Attendee{}, &repository.Category{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Successfully connected to database and migrated schema")

	// Initialize repositories
	repos := repository.NewModels(db)

	// Initialize handlers
	handler := handlers.NewHandler(repos, jwtSecret)

	// Setup router
	router := routers.SetupRouter(handler, jwtSecret, repos.Users)

	// Configure server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start server
	log.Printf("Starting server on port %d", port)
	log.Printf("Swagger UI:   http://localhost:%d/swagger/index.html", port)
	log.Printf("Health Check: http://localhost:%d/health", port)
	log.Println("Database:     Connected")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

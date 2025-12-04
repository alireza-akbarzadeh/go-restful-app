package main

import (
	"database/sql"
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
	_ "modernc.org/sqlite"
)

// @title Go Gin REST API
// @version 1.0
// @description This is a REST API server for managing events and attendees.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	port := config.GetEnvInt("PORT", 8080)
	jwtSecret := config.GetEnvString("JWT_SECRET", "some-secret-123456")
	dbPath := config.GetEnvString("DATABASE_PATH", "./data.db")

	// Initialize database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Successfully connected to database")

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
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

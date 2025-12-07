package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/alireza-akbarzadeh/ginflow/docs"
	"github.com/alireza-akbarzadeh/ginflow/internal/app"
	"github.com/alireza-akbarzadeh/ginflow/internal/console"
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
	port := flag.Int("port", 0, "Server port (overrides PORT env variable)")
	flag.Parse()

	var opts []app.Option
	if *port > 0 {
		opts = append(opts, app.WithPort(*port))
	}

	application, err := app.New(opts...)
	if err != nil {
		c := console.New()
		c.Error("❌", fmt.Sprintf("Failed to initialize application: %v", err))
		os.Exit(1)
	}
	defer application.Shutdown()

	if err := application.Run(); err != nil {
		c := console.New()
		c.Error("❌", fmt.Sprintf("Application error: %v", err))
		os.Exit(1)
	}
}

package main

import (
	_ "github.com/alireza-akbarzadeh/ginflow/docs"
	_ "github.com/joho/godotenv/autoload"

	cmd "github.com/alireza-akbarzadeh/ginflow/cmd/cli"
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
	cmd.Execute()
}

package main

import (
	"log"
	"os"

	"github.com/alireza-akbarzadeh/ginflow/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Drop all tables in the correct order (due to foreign keys)
	err = db.Migrator().DropTable(
		&models.Comment{},
		&models.Attendee{},
		&models.Profile{},
		&models.Event{},
		&models.Category{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("Failed to drop tables:", err)
	}

	log.Println("All tables dropped successfully!")

	// Auto-migrate to recreate tables with integer IDs
	err = db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Event{},
		&models.Profile{},
		&models.Attendee{},
		&models.Comment{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database reset and migrated successfully!")
}

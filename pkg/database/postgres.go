package database

import (
	"fmt"
	"log"

	"github.com/alireza-akbarzadeh/ginflow/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect initializes the database connection
func Connect(dbUrl string) (*gorm.DB, error) {
	if dbUrl == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	dialector := postgres.Open(dbUrl)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Migrate performs auto-migration of database schemas
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	err := db.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Attendee{},
		&models.Category{},
		&models.Comment{},
		&models.Profile{},
		&models.Product{},
	)
	if err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}
	log.Println("Database migration completed successfully")
	return nil
}

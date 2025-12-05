package database

import (
	"fmt"
	"log"
	"time"

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

	var db *gorm.DB
	var err error
	dialector := postgres.Open(dbUrl)

	// Retry connection logic (useful for serverless DBs like Neon that might be sleeping)
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(dialector, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in 2 seconds...", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

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
		&models.BasketItem{},
		&models.Basket{},
	)
	if err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}
	log.Println("Database migration completed successfully")
	return nil
}

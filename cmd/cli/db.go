package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alireza-akbarzadeh/ginflow/internal/console"
	"github.com/alireza-akbarzadeh/ginflow/internal/database"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	forceFlag bool
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
	Long:  `Manage database operations like reset, drop tables, and more.`,
}

// dbResetCmd resets the database (drop all tables and re-migrate)
var dbResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Drop all tables and re-migrate",
	Long:  `Drop all database tables and run migrations to recreate them. This is destructive!`,
	Run: func(cmd *cobra.Command, args []string) {
		c := console.New()

		if !forceFlag && !confirmAction(c, "This will DROP ALL TABLES and reset the database") {
			c.Warning("⚠️", "Operation cancelled.")
			return
		}

		db, err := connectDB(c)
		if err != nil {
			return
		}

		// Drop all tables
		c.Info("🗑️", "Dropping all tables...")
		if err := dropAllTables(db); err != nil {
			c.Error("❌", fmt.Sprintf("Failed to drop tables: %v", err))
			os.Exit(1)
		}
		c.Success("✅", "All tables dropped successfully!")

		// Re-migrate
		c.Info("📦", "Running migrations...")
		if err := database.Migrate(db); err != nil {
			c.Error("❌", fmt.Sprintf("Failed to migrate: %v", err))
			os.Exit(1)
		}
		c.Success("✅", "Database reset and migrated successfully!")
	},
}

// dbDropCmd drops all database tables
var dbDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop all database tables",
	Long:  `Drop all database tables. This is destructive and cannot be undone!`,
	Run: func(cmd *cobra.Command, args []string) {
		c := console.New()

		if !forceFlag && !confirmAction(c, "This will DROP ALL TABLES permanently") {
			c.Warning("⚠️", "Operation cancelled.")
			return
		}

		db, err := connectDB(c)
		if err != nil {
			return
		}

		c.Info("🗑️", "Dropping all tables...")
		if err := dropAllTables(db); err != nil {
			c.Error("❌", fmt.Sprintf("Failed to drop tables: %v", err))
			os.Exit(1)
		}
		c.Success("✅", "All tables dropped successfully!")
	},
}

// dbStatusCmd shows database connection status
var dbStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check database connection status",
	Long:  `Verify the database connection and show basic information.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := console.New()

		db, err := connectDB(c)
		if err != nil {
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			c.Error("❌", fmt.Sprintf("Failed to get database instance: %v", err))
			os.Exit(1)
		}

		// Ping the database
		if err := sqlDB.Ping(); err != nil {
			c.Error("❌", fmt.Sprintf("Database ping failed: %v", err))
			os.Exit(1)
		}

		stats := sqlDB.Stats()
		c.Success("✅", "Database connection is healthy!")
		fmt.Println()
		fmt.Println("Connection Statistics:")
		fmt.Println("────────────────────────────────────")
		fmt.Printf("  Open Connections:    %d\n", stats.OpenConnections)
		fmt.Printf("  In Use:              %d\n", stats.InUse)
		fmt.Printf("  Idle:                %d\n", stats.Idle)
		fmt.Println("────────────────────────────────────")
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbResetCmd)
	dbCmd.AddCommand(dbDropCmd)
	dbCmd.AddCommand(dbStatusCmd)

	// Add --force flag to skip confirmation
	dbResetCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Skip confirmation prompt")
	dbDropCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Skip confirmation prompt")
}

func connectDB(c *console.Console) (*gorm.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		c.Error("❌", "DATABASE_URL environment variable is required")
		os.Exit(1)
	}

	c.Info("🗄️", "Connecting to database...")
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		c.Error("❌", fmt.Sprintf("Failed to connect to database: %v", err))
		os.Exit(1)
	}
	c.Success("✅", "Database connected")
	return db, nil
}

func confirmAction(c *console.Console, message string) bool {
	c.Warning("⚠️", fmt.Sprintf("WARNING: %s. Are you sure? (y/N): ", message))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

func dropAllTables(db *gorm.DB) error {
	// Drop tables in correct order (due to foreign key constraints)
	return db.Migrator().DropTable(
		&models.Comment{},
		&models.Attendee{},
		&models.Profile{},
		&models.Event{},
		&models.BasketItem{},
		&models.Basket{},
		&models.Product{},
		&models.Category{},
		&models.User{},
	)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/alireza-akbarzadeh/ginflow/internal/console"
	"github.com/alireza-akbarzadeh/ginflow/internal/database"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations to create or update tables.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := console.New()

		c.Line()
		c.Info("🔄", "Running database migrations...")
		c.Line()

		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			c.Error("❌", "DATABASE_URL environment variable is not set")
			return fmt.Errorf("DATABASE_URL is required")
		}

		c.Info("🗄️ ", "Connecting to database...")
		db, err := database.Connect(databaseURL)
		if err != nil {
			c.Error("❌", fmt.Sprintf("Failed to connect to database: %v", err))
			return err
		}

		sqlDB, err := db.DB()
		if err != nil {
			c.Error("❌", fmt.Sprintf("Failed to get SQL DB: %v", err))
			return err
		}
		defer sqlDB.Close()

		c.Success("✅", "Database connected")

		c.Info("📦", "Running auto-migrations...")
		if err := database.Migrate(db); err != nil {
			c.Error("❌", fmt.Sprintf("Migration failed: %v", err))
			return err
		}

		c.Line()
		c.Success("✅", "Migrations completed successfully!")
		c.Line()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

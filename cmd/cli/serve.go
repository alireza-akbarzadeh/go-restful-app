package cmd

import (
	"os"

	"github.com/alireza-akbarzadeh/ginflow/internal/app"
	"github.com/alireza-akbarzadeh/ginflow/internal/config"
	"github.com/spf13/cobra"
)

var (
	servePort int
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  `Start the GinFlow HTTP server with all API endpoints.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get configuration from environment
		port := servePort
		if port == 0 {
			port = config.GetEnvInt("PORT", 8080)
		}

		jwtSecret := config.GetEnvString("JWT_SECRET", "")
		if jwtSecret == "" {
			cmd.PrintErrln("Warning: JWT_SECRET not set, using default (not secure for production)")
			jwtSecret = "default-secret-change-me"
		}

		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			return cmd.Help()
		}

		// Create and run the application
		application, err := app.New(
			app.WithPort(port),
			app.WithJWTSecret(jwtSecret),
			app.WithDatabaseURL(databaseURL),
		)
		if err != nil {
			return err
		}
		defer application.Shutdown()

		return application.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Flags for serve command
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 0, "Port to run the server on (default: $PORT or 8080)")
}

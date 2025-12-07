package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const banner = `
   _____ _____ _   _ ______ _      ______          __
  / ____|_   _| \ | |  ____| |    / __ \ \        / /
 | |  __  | | |  \| | |__  | |   | |  | \ \  /\  / / 
 | | |_ | | | | . ` + "`" + ` |  __| | |   | |  | |\ \/  \/ /  
 | |__| |_| |_| |\  | |    | |___| |__| | \  /\  /   
  \_____|_____|_| \_|_|    |______\____/   \/  \/    
`

// Version information (can be set at build time)
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "ginflow",
	Short: "GinFlow - A modern RESTful API framework",
	Long: fmt.Sprintf("\033[36m%s\033[0m\n%s", banner,
		`GinFlow is a production-ready RESTful API framework built with Go and Gin.

Features:
  • RESTful API with CRUD operations
  • JWT Authentication
  • PostgreSQL with GORM
  • Swagger documentation
  • Rate limiting & CORS
  • Graceful shutdown`),
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

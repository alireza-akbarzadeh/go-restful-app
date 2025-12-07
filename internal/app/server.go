package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Run starts the HTTP server and blocks until shutdown signal is received
func (a *App) Run() error {
	// Channel to capture server errors
	serverErr := make(chan error, 1)

	// Start server in goroutine
	go func() {
		a.printStartupBanner()

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case <-quit:
		a.console.Line()
		a.console.Warning("⚠️ ", "Shutdown signal received...")
	}

	return a.gracefulShutdown()
}

// printStartupBanner prints the server startup information
func (a *App) printStartupBanner() {
	// ASCII art for GINFLOW
	banner := `
   _____ _____ _   _ ______ _      ______          __
  / ____|_   _| \ | |  ____| |    / __ \ \        / /
 | |  __  | | |  \| | |__  | |   | |  | \ \  /\  / / 
 | | |_ | | | | . ` + "`" + ` |  __| | |   | |  | |\ \/  \/ /  
 | |__| |_| |_| |\  | |    | |___| |__| | \  /\  /   
  \_____|_____|_| \_|_|    |______\____/   \/  \/    
                                                     `

	fmt.Println("\033[36m" + banner + "\033[0m") // Cyan color

	a.console.Line()
	a.console.Divider()
	a.console.Success("🚀", fmt.Sprintf("Server running on port %d", a.config.Port))
	a.console.Divider()
	a.console.Line()
	a.console.URL("📚", "Swagger UI", fmt.Sprintf("http://localhost:%d/swagger/index.html", a.config.Port))
	a.console.URL("🏥", "Health Check", fmt.Sprintf("http://localhost:%d/health", a.config.Port))
	a.console.URL("🌐", "API Base", fmt.Sprintf("http://localhost:%d/api/v1", a.config.Port))
	a.console.Line()
	a.console.Info("💡", "Press Ctrl+C to stop the server")
	a.console.Line()
}

// gracefulShutdown performs a graceful shutdown of all components
func (a *App) gracefulShutdown() error {
	a.console.Info("🛑", "Gracefully shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), a.config.ShutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	a.console.Line()
	a.console.Success("👋", "Server stopped gracefully. Goodbye!")
	a.console.Line()

	return nil
}

// Shutdown closes all resources (database connections, etc.)
func (a *App) Shutdown() {
	if a.sqlDB != nil {
		a.console.Info("🔌", "Closing database connection...")
		if err := a.sqlDB.Close(); err != nil {
			a.console.Error("❌", fmt.Sprintf("Error closing database: %v", err))
		} else {
			a.console.Success("✅", "Database connection closed")
		}
	}
}

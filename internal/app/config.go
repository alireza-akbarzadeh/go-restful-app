package app

import (
	"time"

	"github.com/alireza-akbarzadeh/ginflow/internal/config"
)

// Config holds all application configuration
type Config struct {
	Port            int
	JWTSecret       string
	DatabaseURL     string
	IdleTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DefaultConfig returns the default configuration loaded from environment
func DefaultConfig() *Config {
	return &Config{
		Port:            config.GetEnvInt("PORT", 8080),
		JWTSecret:       config.GetEnvString("JWT_SECRET", "some-secret-123456"),
		DatabaseURL:     config.GetEnvString("DATABASE_URL", ""),
		IdleTimeout:     time.Minute,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    30 * time.Second,
		ShutdownTimeout: 5 * time.Second,
	}
}

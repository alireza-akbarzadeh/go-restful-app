package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration values
type Config struct {
	// Server Configuration
	Port         string        `env:"PORT" default:"8080"`
	Host         string        `env:"HOST" default:"localhost"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" default:"10s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" default:"10s"`

	// Database Configuration
	DatabaseURL     string        `env:"DATABASE_URL" required:"true"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" default:"5m"`

	// JWT Configuration
	JWTSecret     string        `env:"JWT_SECRET" required:"true"`
	JWTExpiration time.Duration `env:"JWT_EXPIRATION" default:"24h"`

	// Redis Configuration (for future caching)
	RedisURL string `env:"REDIS_URL" default:""`

	// Application Configuration
	Environment string `env:"ENVIRONMENT" default:"development"`
	LogLevel    string `env:"LOG_LEVEL" default:"info"`

	// Rate Limiting
	RateLimitPerMinute int `env:"RATE_LIMIT_PER_MINUTE" default:"100"`
}

// Load loads configuration from environment variables
func Load() *Config {
	config := &Config{
		Port:               getEnv("PORT", "8080"),
		Host:               getEnv("HOST", "localhost"),
		ReadTimeout:        getDurationEnv("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:       getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		MaxOpenConns:       getIntEnv("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:       getIntEnv("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime:    getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpiration:      getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		RedisURL:           getEnv("REDIS_URL", ""),
		Environment:        getEnv("ENVIRONMENT", "development"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		RateLimitPerMinute: getIntEnv("RATE_LIMIT_PER_MINUTE", 100),
	}

	return config
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetAddress returns the full address (host:port)
func (c *Config) GetAddress() string {
	return c.Host + ":" + c.Port
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

package logging

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// ContextKey is the type for context keys
type ContextKey string

const (
	// RequestIDKey is the context key for request ID
	RequestIDKey ContextKey = "request_id"
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
)

var (
	// Logger is the global logger instance
	Logger *slog.Logger
)

// init automatically initializes the logger when the package is imported
func init() {
	InitLogger()
}

// InitLogger initializes the global logger
func InitLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if gin.Mode() == gin.DebugMode {
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context) context.Context {
	requestID := generateRequestID()
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetRequestID gets the request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// GetUserID gets the user ID from context
func GetUserID(ctx context.Context) int {
	if id, ok := ctx.Value(UserIDKey).(int); ok {
		return id
	}
	return 0
}

// LogWithContext creates a logger with context values
func LogWithContext(ctx context.Context) *slog.Logger {
	logger := Logger

	// If Logger is not initialized, use the default slog logger
	if logger == nil {
		logger = slog.Default()
	}

	if requestID := GetRequestID(ctx); requestID != "" {
		logger = logger.With("request_id", requestID)
	}

	if userID := GetUserID(ctx); userID > 0 {
		logger = logger.With("user_id", userID)
	}

	return logger
}

// Info logs an info message with context
func Info(ctx context.Context, msg string, args ...any) {
	LogWithContext(ctx).Info(msg, args...)
}

// Error logs an error message with context
func Error(ctx context.Context, msg string, err error, args ...any) {
	logger := LogWithContext(ctx)
	if err != nil {
		logger = logger.With("error", err.Error())
	}
	logger.Error(msg, args...)
}

// Debug logs a debug message with context
func Debug(ctx context.Context, msg string, args ...any) {
	LogWithContext(ctx).Debug(msg, args...)
}

// Warn logs a warning message with context
func Warn(ctx context.Context, msg string, args ...any) {
	LogWithContext(ctx).Warn(msg, args...)
}

// RequestLoggerMiddleware adds request ID to context and logs requests
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Add request ID to context
		ctx := WithRequestID(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Log request completion
		duration := time.Since(start)
		LogWithContext(ctx).Info("request completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"ip", c.ClientIP(),
		)
	}
}

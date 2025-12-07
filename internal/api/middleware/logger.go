package middleware

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/alireza-akbarzadeh/ginflow/internal/config"
	"github.com/gin-gonic/gin"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorWhite  = "\033[97m"
	colorBold   = "\033[1m"
)

// Logger creates a structured logging middleware
// In development mode (GIN_MODE != release), it uses pretty console output
// In production mode, it uses JSON structured logging
func Logger() gin.HandlerFunc {
	isDev := config.GetEnvString("GIN_MODE", "debug") != "release"

	if isDev {
		return devLogger()
	}
	return prodLogger()
}

// devLogger returns a pretty console logger for development
func devLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate response time
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		// Format latency
		latencyStr := formatLatency(latency)

		// Get colors based on status and method
		statusColor := getStatusColor(statusCode)
		methodColor := getMethodColor(c.Request.Method)

		// Print formatted log
		fmt.Printf("%s %s%-7s%s %s%s%s %s%3d%s %s%s%s %s\n",
			colorGray+time.Now().Format("15:04:05")+colorReset,
			methodColor, c.Request.Method, colorReset,
			colorWhite, c.Request.RequestURI, colorReset,
			statusColor, statusCode, colorReset,
			colorGray, latencyStr, colorReset,
			getStatusEmoji(statusCode),
		)

		// Print errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				fmt.Printf("         %s⚠ Error: %s%s\n", colorRed, err.Error(), colorReset)
			}
		}
	}
}

// prodLogger returns a JSON structured logger for production
func prodLogger() gin.HandlerFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate response time
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		// Log attributes
		attrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.RequestURI),
			slog.String("ip", c.ClientIP()),
			slog.Int("status", statusCode),
			slog.Duration("latency", latency),
		}

		// Log based on status code
		if statusCode >= 500 {
			logger.LogAttrs(c.Request.Context(), slog.LevelError, "Server Error", attrs...)
		} else if statusCode >= 400 {
			logger.LogAttrs(c.Request.Context(), slog.LevelWarn, "Client Error", attrs...)
		} else {
			logger.LogAttrs(c.Request.Context(), slog.LevelInfo, "Request Processed", attrs...)
		}
	}
}

// formatLatency formats the latency duration to a human readable string
func formatLatency(d time.Duration) string {
	switch {
	case d < time.Microsecond:
		return fmt.Sprintf("%dns", d.Nanoseconds())
	case d < time.Millisecond:
		return fmt.Sprintf("%.2fµs", float64(d.Nanoseconds())/1000)
	case d < time.Second:
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000)
	default:
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

// getStatusColor returns ANSI color code based on HTTP status code
func getStatusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return colorGreen + colorBold
	case code >= 300 && code < 400:
		return colorCyan + colorBold
	case code >= 400 && code < 500:
		return colorYellow + colorBold
	default:
		return colorRed + colorBold
	}
}

// getMethodColor returns ANSI color code based on HTTP method
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return colorBlue + colorBold
	case "POST":
		return colorGreen + colorBold
	case "PUT":
		return colorYellow + colorBold
	case "DELETE":
		return colorRed + colorBold
	case "PATCH":
		return colorCyan + colorBold
	default:
		return colorWhite + colorBold
	}
}

// getStatusEmoji returns an emoji based on status code
func getStatusEmoji(code int) string {
	switch {
	case code >= 200 && code < 300:
		return ""
	case code >= 300 && code < 400:
		return "↪️"
	case code == 401:
		return "🔒"
	case code == 403:
		return "🚫"
	case code == 404:
		return "❓"
	case code >= 400 && code < 500:
		return "⚠️"
	default:
		return "❌"
	}
}

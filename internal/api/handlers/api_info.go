package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// APIInfo represents the API information response
type APIInfo struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Docs        string            `json:"docs"`
	Health      string            `json:"health"`
	Endpoints   map[string]string `json:"endpoints"`
	Server      ServerInfo        `json:"server"`
}

// ServerInfo contains server-related information
type ServerInfo struct {
	Timestamp string `json:"timestamp"`
	GoVersion string `json:"go_version"`
	Uptime    string `json:"uptime,omitempty"`
}

var startTime = time.Now()

// GetAPIInfo returns information about the API
// @Summary API Information
// @Description Returns metadata about the GinFlow API including available endpoints and version
// @Tags API
// @Produce json
// @Success 200 {object} APIInfo
// @Router /api/v1 [get]
func (h *Handler) GetAPIInfo(c *gin.Context) {
	baseURL := getBaseURL(c)

	info := APIInfo{
		Name:        "GinFlow API",
		Version:     "v1",
		Description: "A modern RESTful API built with Go and Gin framework",
		Docs:        baseURL + "/swagger/index.html",
		Health:      baseURL + "/health",
		Endpoints: map[string]string{
			"auth":       baseURL + "/api/v1/auth",
			"events":     baseURL + "/api/v1/events",
			"categories": baseURL + "/api/v1/categories",
			"products":   baseURL + "/api/v1/products",
			"attendees":  baseURL + "/api/v1/attendees",
			"profile":    baseURL + "/api/v1/profile",
			"users":      baseURL + "/api/v1/users",
			"basket":     baseURL + "/api/v1/basket",
		},
		Server: ServerInfo{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			GoVersion: runtime.Version(),
			Uptime:    time.Since(startTime).Round(time.Second).String(),
		},
	}

	c.JSON(http.StatusOK, info)
}

// getBaseURL extracts the base URL from the request
func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

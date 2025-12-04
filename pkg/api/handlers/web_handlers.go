package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShowLandingPage renders the HTML landing page
func (h *Handler) ShowLandingPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// ShowHealthPage renders the HTML health status page
func (h *Handler) ShowHealthPage(c *gin.Context) {
	// Check if the client accepts HTML
	if c.GetHeader("Accept") == "application/json" {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}
	c.HTML(http.StatusOK, "health.html", nil)
}

// ShowDashboardPage renders the HTML dashboard page
func (h *Handler) ShowDashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", nil)
}

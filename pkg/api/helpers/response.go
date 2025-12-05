package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// RespondWithError sends an error response
func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	})
}

// RespondWithSuccess sends a success response
func RespondWithSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}

// RespondWithPagination sends a standardized pagination response
func RespondWithPagination(c *gin.Context, data interface{}, total int64, page, limit int, meta map[string]interface{}) {
	response := gin.H{
		"data":  data,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	// Add any additional metadata
	for k, v := range meta {
		response[k] = v
	}

	c.JSON(http.StatusOK, response)
}

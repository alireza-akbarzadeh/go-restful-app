package helpers

import (
	"net/http"

	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/pagination"
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
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

// RespondWithAppError handles AppError types and sends appropriate response
func RespondWithAppError(c *gin.Context, err error, fallbackMessage string) {
	if appErr, ok := err.(*appErrors.AppError); ok {
		response := ErrorResponse{
			Error:   http.StatusText(appErr.StatusCode),
			Message: appErr.Message,
		}
		if len(appErr.Details) > 0 {
			response.Details = appErr.Details
		}
		c.JSON(appErr.StatusCode, response)
		return
	}
	RespondWithError(c, http.StatusInternalServerError, fallbackMessage)
}

// HandleError is a convenience function that handles errors with proper type checking
// It returns true if an error was handled (and response sent), false otherwise
func HandleError(c *gin.Context, err error, fallbackMessage string) bool {
	if err == nil {
		return false
	}
	RespondWithAppError(c, err, fallbackMessage)
	return true
}

// HandleNotFound sends a 404 response if the resource is nil
// Returns true if not found (and response sent), false otherwise
func HandleNotFound(c *gin.Context, resource interface{}, resourceName string) bool {
	if resource == nil {
		RespondWithError(c, http.StatusNotFound, resourceName+" not found")
		return true
	}
	return false
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

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{}                    `json:"data"`
	Pagination *pagination.PaginationResponse `json:"pagination"`
}

// RespondWithPaginatedData sends a paginated response using our pagination structure
func RespondWithPaginatedData(c *gin.Context, code int, data interface{}, paginationResp *pagination.PaginationResponse) {
	c.JSON(code, PaginatedResponse{
		Data:       data,
		Pagination: paginationResp,
	})
}

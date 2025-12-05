package helpers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BindJSON binds the request body to the given interface and handles errors
// Returns true if binding was successful, false otherwise
func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return false
	}
	return true
}

// ParseIDParam extracts and parses an integer ID from URL parameters
func ParseIDParam(c *gin.Context, paramName string) (int, error) {
	idStr := c.Param(paramName)
	return strconv.Atoi(idStr)
}

// ParseQueryInt extracts and parses an integer from query parameters
func ParseQueryInt(c *gin.Context, key string, defaultValue int) int {
	valStr := c.Query(key)
	if valStr == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}

	return val
}

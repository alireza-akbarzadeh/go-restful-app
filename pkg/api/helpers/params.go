package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

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

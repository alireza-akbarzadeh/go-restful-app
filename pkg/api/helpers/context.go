package helpers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/ginflow/pkg/models"
	"github.com/gin-gonic/gin"
)

// GetUserFromContext retrieves the authenticated user from gin context
func GetUserFromContext(c *gin.Context) *models.User {
	contextUser, exists := c.Get("user")
	if !exists {
		return nil
	}
	user, ok := contextUser.(*models.User)
	if !ok {
		return nil
	}
	return user
}

// GetAuthenticatedUser retrieves the authenticated user or sends an unauthorized response
// Returns the user and true if found, nil and false otherwise
func GetAuthenticatedUser(c *gin.Context) (*models.User, bool) {
	user := GetUserFromContext(c)
	if user == nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return nil, false
	}
	return user, true
}

// SetUserInContext sets the authenticated user in gin context
func SetUserInContext(c *gin.Context, user *models.User) {
	c.Set("user", user)
}

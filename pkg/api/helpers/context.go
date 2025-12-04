package helpers

import (
	"github.com/alireza-akbarzadeh/restful-app/pkg/repository"
	"github.com/gin-gonic/gin"
)

// GetUserFromContext retrieves the authenticated user from gin context
func GetUserFromContext(c *gin.Context) *repository.User {
	contextUser, exists := c.Get("user")
	if !exists {
		return nil
	}
	user, ok := contextUser.(*repository.User)
	if !ok {
		return nil
	}
	return user
}

// SetUserInContext sets the authenticated user in gin context
func SetUserInContext(c *gin.Context, user *repository.User) {
	c.Set("user", user)
}

package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/alireza-akbarzadeh/restful-app/pkg/api/helpers"
	"github.com/alireza-akbarzadeh/restful-app/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	msgAuthHeaderRequired    = "Authorization header is required"
	msgAuthHeaderFormat      = "Authorization header format must be Bearer {token}"
	msgInvalidToken          = "Invalid token"
	msgInvalidTokenClaims    = "Invalid token claims"
	msgInvalidUserIDInClaims = "Invalid user ID in token claims"
	msgTokenExpired          = "Token has expired"
	msgUserNotFound          = "User not found"
	msgInternalError         = "An internal error occurred"
)

// AuthMiddleware creates a new authentication middleware
func AuthMiddleware(jwtSecret string, userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			helpers.RespondWithError(c, http.StatusUnauthorized, msgAuthHeaderRequired)
			c.Abort()
			return
		}

		parts := strings.SplitN(authorizationHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			helpers.RespondWithError(c, http.StatusUnauthorized, msgAuthHeaderFormat)
			c.Abort()
			return
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Require HMAC signing to prevent algorithm confusion
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			var e *jwt.ValidationError
			if errors.As(err, &e) && e.Errors&jwt.ValidationErrorExpired != 0 {
				helpers.RespondWithError(c, http.StatusUnauthorized, msgTokenExpired)
			} else {
				helpers.RespondWithError(c, http.StatusUnauthorized, msgInvalidToken)
			}
			c.Abort()
			return
		}

		if !token.Valid {
			helpers.RespondWithError(c, http.StatusUnauthorized, msgInvalidToken)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helpers.RespondWithError(c, http.StatusUnauthorized, msgInvalidTokenClaims)
			c.Abort()
			return
		}

		// Support numeric (float64) and string user_id claim types
		var userID int
		switch v := claims["user_id"].(type) {
		case float64:
			userID = int(v)
		case string:
			i, convErr := strconv.Atoi(v)
			if convErr != nil {
				helpers.RespondWithError(c, http.StatusUnauthorized, msgInvalidUserIDInClaims)
				c.Abort()
				return
			}
			userID = i
		default:
			helpers.RespondWithError(c, http.StatusUnauthorized, msgInvalidUserIDInClaims)
			c.Abort()
			return
		}

		user, err := userRepo.Get(userID)
		if err != nil {
			helpers.RespondWithError(c, http.StatusInternalServerError, msgInternalError)
			c.Abort()
			return
		}
		if user == nil {
			helpers.RespondWithError(c, http.StatusUnauthorized, msgUserNotFound)
			c.Abort()
			return
		}

		helpers.SetUserInContext(c, user)
		c.Next()
	}
}

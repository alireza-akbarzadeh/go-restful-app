package handlers

import (
	"net/http"
	"time"

	"github.com/alireza-akbarzadeh/restful-app/pkg/api/helpers"
	"github.com/alireza-akbarzadeh/restful-app/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=3"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string           `json:"token"`
	User  *repository.User `json:"user"`
}

// Login handles user authentication
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest   true  "Login credentials"
// @Success      200          {object}  LoginResponse
// @Failure      400          {object}  helpers.ErrorResponse
// @Failure      401          {object}  helpers.ErrorResponse
// @Failure      500          {object}  helpers.ErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get user by email
	user, err := h.Repos.Users.GetByEmail(req.Email)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if user == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Don't expose password in response
	user.Password = ""

	c.JSON(http.StatusOK, LoginResponse{
		Token: tokenString,
		User:  user,
	})
}

// Register handles user registration
// @Summary      User registration
// @Description  Register a new user account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      RegisterRequest  true  "User registration details"
// @Success      201   {object}  repository.User
// @Failure      400   {object}  helpers.ErrorResponse
// @Failure      500   {object}  helpers.ErrorResponse
// @Router       /api/v1/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if user already exists
	existingUser, err := h.Repos.Users.GetByEmail(req.Email)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to check existing user")
		return
	}
	if existingUser != nil {
		helpers.RespondWithError(c, http.StatusConflict, "User with this email already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to register user")
		return
	}

	// Create user
	user := &repository.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	createdUser, err := h.Repos.Users.Insert(user)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to register user")
		return
	}

	// Don't expose password in response
	createdUser.Password = ""

	c.JSON(http.StatusCreated, createdUser)
}

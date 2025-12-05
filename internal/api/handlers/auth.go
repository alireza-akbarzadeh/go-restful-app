package handlers

import (
	"net/http"
	"time"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/helpers"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
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
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// UpdatePasswordRequest represents the password update request payload
type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
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
	user, err := h.Repos.Users.GetByEmail(c.Request.Context(), req.Email)
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

	// Update last login time
	_ = h.Repos.Users.UpdateLastLogin(c.Request.Context(), user.ID)

	// Don't expose password in response
	user.Password = ""

	c.JSON(http.StatusOK, LoginResponse{
		Token: tokenString,
		User:  user,
	})
}

// Logout handles user logout
// @Summary      User logout
// @Description  Invalidate user session (client-side)
// @Tags         Authentication
// @Success      200          {object}  map[string]string
// @Router       /api/v1/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// Register handles user registration
// @Summary      User registration
// @Description  Register a new user account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      RegisterRequest  true  "User registration details"
// @Success      201   {object}  models.User
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
	existingUser, err := h.Repos.Users.GetByEmail(c.Request.Context(), req.Email)
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
	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	createdUser, err := h.Repos.Users.Insert(c.Request.Context(), user)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to register user")
		return
	}

	// Don't expose password in response
	createdUser.Password = ""

	c.JSON(http.StatusCreated, createdUser)
}

// UpdatePassword handles password updates
// @Summary      Update password
// @Description  Update the authenticated user's password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body UpdatePasswordRequest true "Password update details"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  helpers.ErrorResponse
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/auth/password [put]
func (h *Handler) UpdatePassword(c *gin.Context) {
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get authenticated user
	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get fresh user data to verify old password
	dbUser, err := h.Repos.Users.Get(c.Request.Context(), user.ID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}
	if dbUser == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(req.OldPassword))
	if err != nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Invalid old password")
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update password
	if err := h.Repos.Users.UpdatePassword(c.Request.Context(), user.ID, string(hashedPassword)); err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to update password")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

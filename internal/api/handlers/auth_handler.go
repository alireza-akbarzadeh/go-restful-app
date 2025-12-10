package handlers

import (
	"net/http"
	"time"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/helpers"
	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
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
	ctx := c.Request.Context()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	logging.Debug(ctx, "login attempt", "email", req.Email)

	// Get user by email
	user, err := h.Repos.Users.GetByEmail(ctx, req.Email)
	if err != nil {
		// Check if it's a not found error - that's fine, just invalid credentials
		if appErrors.IsType(err, appErrors.ErrNotFound) {
			helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrUnauthorized, "Invalid email or password"), "")
			return
		}
		helpers.HandleError(c, err, "Something went wrong")
		return
	}

	// Check if user was found
	if user == nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrUnauthorized, "Invalid email or password"), "")
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrUnauthorized, "Invalid email or password"), "")
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		logging.Error(ctx, "failed to generate token", err, "user_id", user.ID)
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Update last login time
	_ = h.Repos.Users.UpdateLastLogin(ctx, user.ID)

	// Don't expose password in response
	user.Password = ""

	logging.Info(ctx, "user logged in successfully", "user_id", user.ID, "email", user.Email)
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
	ctx := c.Request.Context()

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	logging.Debug(ctx, "registration attempt", "email", req.Email)

	// Check if user already exists
	existingUser, err := h.Repos.Users.GetByEmail(ctx, req.Email)
	if err != nil && !appErrors.IsType(err, appErrors.ErrNotFound) {
		helpers.HandleError(c, err, "Failed to check existing user")
		return
	}
	if existingUser != nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrAlreadyExists, "User with this email already exists"), "")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logging.Error(ctx, "failed to hash password", err)
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to register user")
		return
	}

	// Create user
	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	createdUser, err := h.Repos.Users.Insert(ctx, user)
	if helpers.HandleError(c, err, "Failed to register user") {
		return
	}

	// Don't expose password in response
	createdUser.Password = ""

	logging.Info(ctx, "user registered successfully", "user_id", createdUser.ID, "email", createdUser.Email)
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
	ctx := c.Request.Context()

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get authenticated user
	user, ok := helpers.GetAuthenticatedUser(c)
	if !ok {
		return
	}

	logging.Debug(ctx, "password update attempt", "user_id", user.ID)

	// Get fresh user data to verify old password
	dbUser, err := h.Repos.Users.Get(ctx, user.ID)
	if helpers.HandleError(c, err, "Failed to retrieve user") {
		return
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(req.OldPassword))
	if err != nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrUnauthorized, "Invalid old password"), "")
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logging.Error(ctx, "failed to hash new password", err, "user_id", user.ID)
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update password
	if err := h.Repos.Users.UpdatePassword(ctx, user.ID, string(hashedPassword)); err != nil {
		helpers.HandleError(c, err, "Failed to update password")
		return
	}

	logging.Info(ctx, "password updated successfully", "user_id", user.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

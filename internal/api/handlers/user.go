package handlers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/helpers"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/gin-gonic/gin"
)

// GetAllUsers retrieves all users
// @Summary      Get all users
// @Description  Get a list of all registered users
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.Repos.Users.GetAll(c.Request.Context())
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser updates a user's profile
// @Summary      Update user profile
// @Description  Update user details (Name, Email)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      int              true  "User ID"
// @Param        user  body      models.User  true  "User object"
// @Success      200   {object}  models.User
// @Failure      400   {object}  helpers.ErrorResponse
// @Failure      401   {object}  helpers.ErrorResponse
// @Failure      403   {object}  helpers.ErrorResponse
// @Failure      404   {object}  helpers.ErrorResponse
// @Failure      500   {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get authenticated user
	authUser, ok := helpers.GetAuthenticatedUser(c)
	if !ok {
		return
	}

	// Allow users to update only their own profile (unless we add admin role later)
	if authUser.ID != id {
		helpers.RespondWithError(c, http.StatusForbidden, "You can only update your own profile")
		return
	}

	// Get existing user
	existingUser, err := h.Repos.Users.Get(c.Request.Context(), id)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}
	if existingUser == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Bind new data
	var updateData models.User
	if !helpers.BindJSON(c, &updateData) {
		return
	}

	// Update fields (preserve ID and Password)
	existingUser.Name = updateData.Name
	existingUser.Email = updateData.Email
	// Note: Password update is handled by a separate endpoint

	if err := h.Repos.Users.Update(c.Request.Context(), existingUser); err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	c.JSON(http.StatusOK, existingUser)
}

// DeleteUser deletes a user
// @Summary      Delete user
// @Description  Delete a user account
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      204
// @Failure      400  {object}  helpers.ErrorResponse
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      403  {object}  helpers.ErrorResponse
// @Failure      404  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get authenticated user
	authUser := helpers.GetUserFromContext(c)
	if authUser == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Allow users to delete only their own profile
	if authUser.ID != id {
		helpers.RespondWithError(c, http.StatusForbidden, "You can only delete your own profile")
		return
	}

	if err := h.Repos.Users.Delete(c.Request.Context(), id); err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	c.Status(http.StatusNoContent)
}

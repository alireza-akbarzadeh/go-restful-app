package handlers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/helpers"
	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/query"
	"github.com/gin-gonic/gin"
)

// GetAllUsers retrieves all users with pagination, filtering, sorting, and search
// @Summary      Get all users
// @Description  Get a paginated list of all registered users with filtering, sorting, and search
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        page        query     int     false  "Page number (default: 1)"
// @Param        page_size   query     int     false  "Page size (default: 20, max: 100)"
// @Param        type        query     string  false  "Pagination type: 'offset' or 'cursor' (default: offset)"
// @Param        cursor      query     string  false  "Cursor for cursor-based pagination"
// @Param        sort        query     string  false  "Sort fields (e.g., '-created_at,name:asc')"
// @Param        search      query     string  false  "Search term for name, email"
// @Param        name[like]  query     string  false  "Filter by name (partial match)"
// @Param        email[like] query     string  false  "Filter by email (partial match)"
// @Param        created_at[gte] query string false  "Filter by minimum created date"
// @Param        created_at[lte] query string false  "Filter by maximum created date"
// @Success      200  {object}  query.PaginatedList{data=[]models.User}
// @Failure      400  {object}  helpers.ErrorResponse
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters from context
	params := query.ParseFromContext(c)

	logging.Debug(ctx, "retrieving all users",
		"page", params.Page,
		"page_size", params.PageSize,
		"search", params.Search,
	)

	users, result, err := h.Repos.Users.GetAll(ctx, params)
	if helpers.HandleError(c, err, "Failed to retrieve users") {
		return
	}

	logging.Debug(ctx, "users retrieved successfully",
		"count", len(users),
		"page", params.Page,
	)

	c.JSON(http.StatusOK, result)
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

	ctx := c.Request.Context()
	logging.Debug(ctx, "updating user", "user_id", id, "auth_user_id", authUser.ID)

	// Get existing user
	existingUser, err := h.Repos.Users.Get(ctx, id)
	if err != nil {
		logging.Error(ctx, "failed to get user for update", err, "user_id", id)
		if appErr, ok := err.(*appErrors.AppError); ok {
			helpers.RespondWithError(c, appErr.StatusCode, appErr.Message)
		} else {
			helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve user")
		}
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

	if err := h.Repos.Users.Update(ctx, existingUser); err != nil {
		logging.Error(ctx, "failed to update user", err, "user_id", id)
		if appErr, ok := err.(*appErrors.AppError); ok {
			helpers.RespondWithError(c, appErr.StatusCode, appErr.Message)
		} else {
			helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to update user")
		}
		return
	}

	logging.Info(ctx, "user updated successfully", "user_id", id, "email", existingUser.Email)
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

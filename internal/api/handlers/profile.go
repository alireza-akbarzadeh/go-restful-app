package handlers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/helpers"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/gin-gonic/gin"
)

// GetProfile retrieves the current user's profile
// @Summary      Get user profile
// @Description  Get the authenticated user's profile information
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Profile
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      404  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	// Get Authenticated User
	authUser := helpers.GetUserFromContext(c)
	if authUser == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// Get Profile with user Data
	profile, err := h.Repos.Profiles.GetByUserIDWithUser(c.Request.Context(), authUser.ID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve profile")
		return
	}
	if profile == nil {
		profile = &models.Profile{
			UserID: authUser.ID,
			User:   *authUser,
		}
	}
	c.JSON(http.StatusOK, profile)
}

// CreateProfile creates a new profile for the authenticated user
// @Summary      Create user profile
// @Description  Create a new profile for the authenticated user
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        profile  body      models.Profile  true  "Profile object"
// @Success      201      {object}  models.Profile
// @Failure      400      {object}  helpers.ErrorResponse
// @Failure      401      {object}  helpers.ErrorResponse
// @Failure      409      {object}  helpers.ErrorResponse
// @Failure      500      {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/profile [post]
func (h *Handler) CreateProfile(c *gin.Context) {
	// Get an authenticated user
	authUser := helpers.GetUserFromContext(c)
	if authUser == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// Check if a profile already exists
	existingProfile, err := h.Repos.Profiles.GetByUserID(c.Request.Context(), authUser.ID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve profile")
		return
	}
	if existingProfile != nil {
		helpers.RespondWithError(c, http.StatusConflict, "Profile already exists")
		return
	}
	// Bind Profile Data
	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	// Set user ID and create a profile
	profile.UserID = authUser.ID
	createdProfile, err := h.Repos.Profiles.Insert(c.Request.Context(), &profile)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to create profile")
		return
	}
	// Return profile with user data
	profileWithUser, err := h.Repos.Profiles.GetByUserIDWithUser(c.Request.Context(), authUser.ID)
	if err != nil {
		// If preloading fails, return the created profile
		c.JSON(http.StatusCreated, createdProfile)
		return
	}
	c.JSON(http.StatusCreated, profileWithUser)
}

// UpdateProfile updates the authenticated user's profile
// @Summary      Update user profile
// @Description  Update the authenticated user's profile information
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        profile  body      models.Profile  true  "Profile object"
// @Success      200      {object}  models.Profile
// @Failure      400      {object}  helpers.ErrorResponse
// @Failure      401      {object}  helpers.ErrorResponse
// @Failure      404      {object}  helpers.ErrorResponse
// @Failure      500      {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	authUser := helpers.GetUserFromContext(c)
	if authUser == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// Get an existing profile
	existingProfile, err := h.Repos.Profiles.GetByUserID(c.Request.Context(), authUser.ID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve profile")
		return
	}

	// Bind update data
	var updateData models.Profile
	if err := c.ShouldBindJSON(&updateData); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	// If a profile doesn't exist, create it
	if existingProfile == nil {
		updateData.UserID = authUser.ID
		createdProfile, err := h.Repos.Profiles.Insert(c.Request.Context(), &updateData)
		if err != nil {
			helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to create profile")
			return
		}
		profileWithUser, err := h.Repos.Profiles.GetByUserIDWithUser(c.Request.Context(), authUser.ID)
		if err != nil {
			c.JSON(http.StatusCreated, createdProfile)
			return
		}
		c.JSON(http.StatusCreated, profileWithUser)
		return
	}

	// Update existing profile fields
	existingProfile.Bio = updateData.Bio
	existingProfile.AvatarURL = updateData.AvatarURL
	existingProfile.Phone = updateData.Phone
	existingProfile.DateOfBirth = updateData.DateOfBirth
	existingProfile.Country = updateData.Country
	existingProfile.City = updateData.City
	existingProfile.Timezone = updateData.Timezone
	existingProfile.Website = updateData.Website
	existingProfile.Twitter = updateData.Twitter
	existingProfile.LinkedIn = updateData.LinkedIn
	existingProfile.GitHub = updateData.GitHub
	existingProfile.Theme = updateData.Theme
	existingProfile.Language = updateData.Language
	existingProfile.IsPublic = updateData.IsPublic
	existingProfile.EmailNotifications = updateData.EmailNotifications
	existingProfile.PushNotifications = updateData.PushNotifications
	//save updated
	if err := h.Repos.Profiles.Update(c.Request.Context(), existingProfile); err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to update profile")

	}
	c.JSON(http.StatusOK, existingProfile)
	// Return an updated profile with user data
	profileWithUser, err := h.Repos.Profiles.GetByUserIDWithUser(c.Request.Context(), authUser.ID)
	if err != nil {
		c.JSON(http.StatusOK, existingProfile)
		return
	}
	c.JSON(http.StatusOK, profileWithUser)
}

// DeleteProfile deletes the authenticated user's profile
// @Summary      Delete user profile
// @Description  Delete the authenticated user's profile
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Success      204  {object}  nil
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      404  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/profile [delete]
func (h *Handler) DeleteProfile(c *gin.Context) {
	// Get an authenticated user
	authUser := helpers.GetUserFromContext(c)
	if authUser == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// Check if a profile exists
	existingProfile, err := h.Repos.Profiles.GetByUserID(c.Request.Context(), authUser.ID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve profile")
		return
	}
	if existingProfile == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "Profile not found")
		return
	}
	// Delete profile
	if err := h.Repos.Profiles.DeleteByUserID(c.Request.Context(), authUser.ID); err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to delete profile")
		return
	}

	c.Status(http.StatusNoContent)
}

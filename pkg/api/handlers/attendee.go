package handlers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/restful-app/pkg/api/helpers"
	"github.com/alireza-akbarzadeh/restful-app/pkg/repository"
	"github.com/gin-gonic/gin"
)

// AddAttendee adds an attendee to an event
// @Summary      Add attendee to event
// @Description  Add an attendee to an event (requires authentication and event ownership)
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Param        id      path      int  true  "Event ID"
// @Param        userId  path      int  true  "User ID"
// @Success      201     {object}  repository.Attendee
// @Failure      400     {object}  helpers.ErrorResponse
// @Failure      401     {object}  helpers.ErrorResponse
// @Failure      403     {object}  helpers.ErrorResponse
// @Failure      404     {object}  helpers.ErrorResponse
// @Failure      409     {object}  helpers.ErrorResponse
// @Failure      500     {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/events/{id}/attendees/{userId} [post]
func (h *Handler) AddAttendee(c *gin.Context) {
	eventID, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	userID, err := helpers.ParseIDParam(c, "userId")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get authenticated user
	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check if event exists
	event, err := h.Repos.Events.Get(eventID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve event")
		return
	}
	if event == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "Event not found")
		return
	}

	// Check if user is the event owner
	if event.OwnerID != user.ID {
		helpers.RespondWithError(c, http.StatusForbidden, "You are not authorized to add attendees to this event")
		return
	}

	// Check if user to add exists
	userToAdd, err := h.Repos.Users.Get(userID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}
	if userToAdd == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Check if already attending
	existingAttendee, err := h.Repos.Attendees.GetByEventAndUser(eventID, userID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to check attendee status")
		return
	}
	if existingAttendee != nil {
		helpers.RespondWithError(c, http.StatusConflict, "User is already attending this event")
		return
	}

	// Add attendee
	attendee := &repository.Attendee{
		EventID: eventID,
		UserID:  userID,
	}

	createdAttendee, err := h.Repos.Attendees.Insert(attendee)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to add attendee")
		return
	}

	c.JSON(http.StatusCreated, createdAttendee)
}

// GetAttendees retrieves all attendees for an event
// @Summary      Get attendees for event
// @Description  Get all attendees for a specific event
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Event ID"
// @Success      200 {array}   repository.User
// @Failure      400 {object}  helpers.ErrorResponse
// @Failure      404 {object}  helpers.ErrorResponse
// @Failure      500 {object}  helpers.ErrorResponse
// @Router       /api/v1/events/{id}/attendees [get]
func (h *Handler) GetAttendees(c *gin.Context) {
	eventID, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	// Check if event exists
	event, err := h.Repos.Events.Get(eventID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve event")
		return
	}
	if event == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "Event not found")
		return
	}

	users, err := h.Repos.Attendees.GetAttendeesByEvent(eventID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve attendees")
		return
	}

	c.JSON(http.StatusOK, users)
}

// RemoveAttendee removes an attendee from an event
// @Summary      Remove attendee from event
// @Description  Remove an attendee from an event (requires authentication and event ownership)
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Param        id      path  int  true  "Event ID"
// @Param        userId  path  int  true  "User ID"
// @Success      204
// @Failure      400  {object}  helpers.ErrorResponse
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      403  {object}  helpers.ErrorResponse
// @Failure      404  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/events/{id}/attendees/{userId} [delete]
func (h *Handler) RemoveAttendee(c *gin.Context) {
	eventID, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	userID, err := helpers.ParseIDParam(c, "userId")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get authenticated user
	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check if event exists
	event, err := h.Repos.Events.Get(eventID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve event")
		return
	}
	if event == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "Event not found")
		return
	}

	// Check if user is the event owner
	if event.OwnerID != user.ID {
		helpers.RespondWithError(c, http.StatusForbidden, "You are not authorized to remove attendees from this event")
		return
	}

	// Remove attendee
	if err := h.Repos.Attendees.Delete(userID, eventID); err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to remove attendee")
		return
	}

	c.Status(http.StatusNoContent)
}

// GetEventsByAttendee retrieves all events for a specific attendee
// @Summary      Get events by attendee
// @Description  Get all events for a specific attendee
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Attendee/User ID"
// @Success      200 {array}   repository.Event
// @Failure      400 {object}  helpers.ErrorResponse
// @Failure      500 {object}  helpers.ErrorResponse
// @Router       /api/v1/attendees/{id}/events [get]
func (h *Handler) GetEventsByAttendee(c *gin.Context) {
	userID, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	events, err := h.Repos.Attendees.GetEventsByAttendee(userID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}

	c.JSON(http.StatusOK, events)
}

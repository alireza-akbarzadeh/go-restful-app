package handlers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/helpers"
	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
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
// @Success      201     {object}  models.Attendee
// @Failure      400     {object}  helpers.ErrorResponse
// @Failure      401     {object}  helpers.ErrorResponse
// @Failure      403     {object}  helpers.ErrorResponse
// @Failure      404     {object}  helpers.ErrorResponse
// @Failure      409     {object}  helpers.ErrorResponse
// @Failure      500     {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/events/{id}/attendees/{userId} [post]
func (h *Handler) AddAttendee(c *gin.Context) {
	ctx := c.Request.Context()

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
	authUser, ok := helpers.GetAuthenticatedUser(c)
	if !ok {
		return
	}

	logging.Debug(ctx, "adding attendee to event", "event_id", eventID, "user_id", userID, "auth_user_id", authUser.ID)

	// Check if the event exists
	event, err := h.Repos.Events.Get(ctx, eventID)
	if helpers.HandleError(c, err, "Failed to retrieve event") {
		return
	}

	// Check if user is the event owner
	if event.OwnerID != authUser.ID {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrForbidden, "You are not authorized to add attendees to this event"), "")
		return
	}

	// Check if user to add exists
	_, err = h.Repos.Users.Get(ctx, userID)
	if helpers.HandleError(c, err, "Failed to retrieve user") {
		return
	}

	// Check if already attending
	existingAttendee, err := h.Repos.Attendees.GetByEventAndUser(ctx, eventID, userID)
	if err != nil {
		helpers.HandleError(c, err, "Failed to check attendee status")
		return
	}
	if existingAttendee != nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrAlreadyExists, "User is already attending this event"), "")
		return
	}

	// Add attendee
	attendee := &models.Attendee{
		EventID: eventID,
		UserID:  userID,
	}

	createdAttendee, err := h.Repos.Attendees.Insert(ctx, attendee)
	if helpers.HandleError(c, err, "Failed to add attendee") {
		return
	}

	logging.Info(ctx, "attendee added successfully", "attendee_id", createdAttendee.ID, "event_id", eventID, "user_id", userID)
	c.JSON(http.StatusCreated, createdAttendee)
}

// GetAttendees retrieves all attendees for an event
// @Summary      Get attendees for event
// @Description  Get all attendees for a specific event
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Event ID"
// @Success      200 {array}   models.User
// @Failure      400 {object}  helpers.ErrorResponse
// @Failure      404 {object}  helpers.ErrorResponse
// @Failure      500 {object}  helpers.ErrorResponse
// @Router       /api/v1/events/{id}/attendees [get]
func (h *Handler) GetAttendees(c *gin.Context) {
	ctx := c.Request.Context()

	eventID, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	logging.Debug(ctx, "retrieving attendees for event", "event_id", eventID)

	// Check if the event exists
	_, err = h.Repos.Events.Get(ctx, eventID)
	if helpers.HandleError(c, err, "Failed to retrieve event") {
		return
	}

	users, err := h.Repos.Attendees.GetAttendeesByEvent(ctx, eventID)
	if helpers.HandleError(c, err, "Failed to retrieve attendees") {
		return
	}

	logging.Debug(ctx, "attendees retrieved successfully", "event_id", eventID, "count", len(users))
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
	ctx := c.Request.Context()

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
	authUser, ok := helpers.GetAuthenticatedUser(c)
	if !ok {
		return
	}

	logging.Debug(ctx, "removing attendee from event", "event_id", eventID, "user_id", userID, "auth_user_id", authUser.ID)

	// Check if the event exists
	event, err := h.Repos.Events.Get(ctx, eventID)
	if helpers.HandleError(c, err, "Failed to retrieve event") {
		return
	}

	// Check if user is the event owner
	if event.OwnerID != authUser.ID {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrForbidden, "You are not authorized to remove attendees from this event"), "")
		return
	}

	// Remove attendee
	if err := h.Repos.Attendees.Delete(ctx, userID, eventID); err != nil {
		helpers.HandleError(c, err, "Failed to remove attendee")
		return
	}

	logging.Info(ctx, "attendee removed successfully", "event_id", eventID, "user_id", userID)
	c.Status(http.StatusNoContent)
}

// GetEventsByAttendee retrieves all events for a specific attendee
// @Summary      Get events by attendee
// @Description  Get all events for a specific attendee
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Attendee/User ID"
// @Success      200 {array}   models.Event
// @Failure      400 {object}  helpers.ErrorResponse
// @Failure      500 {object}  helpers.ErrorResponse
// @Router       /api/v1/attendees/{id}/events [get]
func (h *Handler) GetEventsByAttendee(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	logging.Debug(ctx, "retrieving events for attendee", "user_id", userID)

	events, err := h.Repos.Attendees.GetEventsByAttendee(ctx, userID)
	if helpers.HandleError(c, err, "Failed to retrieve events") {
		return
	}

	logging.Debug(ctx, "events retrieved for attendee", "user_id", userID, "count", len(events))
	c.JSON(http.StatusOK, events)
}

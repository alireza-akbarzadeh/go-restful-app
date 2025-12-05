package handlers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/ginflow/pkg/api/helpers"
	"github.com/alireza-akbarzadeh/ginflow/pkg/models"
	"github.com/gin-gonic/gin"
)

// CreateEvent handles event creation
// @Summary      Create a new event
// @Description  Create a new event (requires authentication)
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        event  body      models.Event  true  "Event object"
// @Success      201    {object}  models.Event
// @Failure      400    {object}  helpers.ErrorResponse
// @Failure      401    {object}  helpers.ErrorResponse
// @Failure      500    {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/events [post]
func (h *Handler) CreateEvent(c *gin.Context) {
	var event models.Event
	if !helpers.BindJSON(c, &event) {
		return
	}

	// Get authenticated user
	user, ok := helpers.GetAuthenticatedUser(c)
	if !ok {
		return
	}

	event.OwnerID = user.ID
	createdEvent, err := h.Repos.Events.Insert(&event)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to create event")
		return
	}

	c.JSON(http.StatusCreated, createdEvent)
}

// GetEvent retrieves a single event by ID
// @Summary      Get a single event
// @Description  Get event by ID
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  models.Event
// @Failure      400  {object}  helpers.ErrorResponse
// @Failure      404  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Router       /api/v1/events/{id} [get]
func (h *Handler) GetEvent(c *gin.Context) {
	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	event, err := h.Repos.Events.Get(id)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve event")
		return
	}
	if event == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "Event not found")
		return
	}

	c.JSON(http.StatusOK, event)
}

// GetAllEvents retrieves all events
// @Summary      Get all events
// @Description  Get a list of all events
// @Tags         Events
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Event
// @Failure      500  {object}  helpers.ErrorResponse
// @Router       /api/v1/events [get]
func (h *Handler) GetAllEvents(c *gin.Context) {
	events, err := h.Repos.Events.GetAll()
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}

	c.JSON(http.StatusOK, events)
}

// UpdateEvent updates an existing event
// @Summary      Update an event
// @Description  Update an existing event (requires authentication and ownership)
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id     path      int             true  "Event ID"
// @Param        event  body      models.Event  true  "Event object"
// @Success      200    {object}  models.Event
// @Failure      400    {object}  helpers.ErrorResponse
// @Failure      401    {object}  helpers.ErrorResponse
// @Failure      403    {object}  helpers.ErrorResponse
// @Failure      404    {object}  helpers.ErrorResponse
// @Failure      500    {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/events/{id} [put]
func (h *Handler) UpdateEvent(c *gin.Context) {
	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	// Get authenticated user
	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check if event exists and user is the owner
	existingEvent, err := h.Repos.Events.Get(id)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve event")
		return
	}
	if existingEvent == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "Event not found")
		return
	}
	if existingEvent.OwnerID != user.ID {
		helpers.RespondWithError(c, http.StatusForbidden, "You are not authorized to update this event")
		return
	}

	// Parse updated event data
	var updatedEvent models.Event
	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedEvent.ID = id
	updatedEvent.OwnerID = user.ID

	if err := h.Repos.Events.Update(&updatedEvent); err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to update event")
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

// DeleteEvent deletes an event
// @Summary      Delete an event
// @Description  Delete an existing event (requires authentication and ownership)
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Event ID"
// @Success      204
// @Failure      400  {object}  helpers.ErrorResponse
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      403  {object}  helpers.ErrorResponse
// @Failure      404  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/events/{id} [delete]
func (h *Handler) DeleteEvent(c *gin.Context) {
	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	// Get authenticated user
	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check if event exists and user is the owner
	existingEvent, err := h.Repos.Events.Get(id)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve event")
		return
	}
	if existingEvent == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "Event not found")
		return
	}
	if existingEvent.OwnerID != user.ID {
		helpers.RespondWithError(c, http.StatusForbidden, "You are not authorized to delete this event")
		return
	}

	// Delete all attendees first
	if err := h.Repos.Attendees.DeleteByEvent(id); err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to delete event attendees")
		return
	}

	// Delete event
	if err := h.Repos.Events.Delete(id); err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to delete event")
		return
	}

	c.Status(http.StatusNoContent)
}

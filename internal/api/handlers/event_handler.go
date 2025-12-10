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

	ctx := c.Request.Context()
	logging.Debug(ctx, "creating new event", "name", event.Name, "owner_id", user.ID)

	event.OwnerID = user.ID
	createdEvent, err := h.Repos.Events.Insert(ctx, &event)
	if helpers.HandleError(c, err, "Failed to create event") {
		return
	}

	logging.Info(ctx, "event created successfully", "event_id", createdEvent.ID, "name", createdEvent.Name)
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

	ctx := c.Request.Context()
	logging.Debug(ctx, "retrieving event", "event_id", id)

	event, err := h.Repos.Events.Get(ctx, id)
	if helpers.HandleError(c, err, "Failed to retrieve event") {
		return
	}

	if event == nil {
		helpers.RespondWithAppError(c, appErrors.Newf(appErrors.ErrNotFound, "event with ID %d not found", id), "")
		return
	}

	logging.Debug(ctx, "event retrieved successfully", "event_id", id, "name", event.Name)
	c.JSON(http.StatusOK, event)
}

// GetAllEvents retrieves all events with advanced pagination
// @Summary      Get all events
// @Description  Get a paginated list of all events with filtering, sorting, and search
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        page        query     int     false  "Page number (default: 1)"
// @Param        page_size   query     int     false  "Page size (default: 20, max: 100)"
// @Param        type        query     string  false  "Pagination type: 'offset' or 'cursor' (default: offset)"
// @Param        cursor      query     string  false  "Cursor for cursor-based pagination"
// @Param        sort        query     string  false  "Sort fields (e.g., '-created_at,name:asc')"
// @Param        search      query     string  false  "Search term for name, description, location"
// @Param        name[eq]    query     string  false  "Filter by exact name"
// @Param        name[like]  query     string  false  "Filter by name (partial match)"
// @Param        location[eq] query    string  false  "Filter by location"
// @Param        owner_id[eq] query    int     false  "Filter by owner ID"
// @Success      200  {object}  query.PaginatedList{data=[]models.Event}
// @Failure      400  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Router       /api/v1/events [get]
func (h *Handler) GetAllEvents(c *gin.Context) {
	ctx := c.Request.Context()
	logging.Debug(ctx, "handling GetAllEvents request with advanced pagination")

	// Parse advanced pagination parameters from context
	req := query.ParseFromContext(c)

	events, result, err := h.Repos.Events.ListWithAdvancedPagination(ctx, req)
	if helpers.HandleError(c, err, "Failed to retrieve events") {
		return
	}

	logging.Info(ctx, "events retrieved successfully",
		"count", len(events),
		"page", req.Page,
		"type", req.Type,
	)

	c.JSON(http.StatusOK, result)
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
	ctx := c.Request.Context()

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	// Get authenticated user
	user, ok := helpers.GetAuthenticatedUser(c)
	if !ok {
		return
	}

	logging.Debug(ctx, "updating event", "event_id", id, "user_id", user.ID)

	// Check if event exists and user is the owner
	existingEvent, err := h.Repos.Events.Get(ctx, id)
	if helpers.HandleError(c, err, "Failed to retrieve event") {
		return
	}

	if existingEvent.OwnerID != user.ID {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrForbidden, "You are not authorized to update this event"), "")
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

	if err := h.Repos.Events.Update(ctx, &updatedEvent); err != nil {
		helpers.HandleError(c, err, "Failed to update event")
		return
	}

	logging.Info(ctx, "event updated successfully", "event_id", id, "name", updatedEvent.Name)
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
	ctx := c.Request.Context()

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	// Get authenticated user
	user, ok := helpers.GetAuthenticatedUser(c)
	if !ok {
		return
	}

	logging.Debug(ctx, "deleting event", "event_id", id, "user_id", user.ID)

	// Check if event exists and user is the owner
	existingEvent, err := h.Repos.Events.Get(ctx, id)
	if helpers.HandleError(c, err, "Failed to retrieve event") {
		return
	}

	if existingEvent.OwnerID != user.ID {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrForbidden, "You are not authorized to delete this event"), "")
		return
	}

	// Delete all attendees first
	if err := h.Repos.Attendees.DeleteByEvent(ctx, id); err != nil {
		helpers.HandleError(c, err, "Failed to delete event attendees")
		return
	}

	// Delete event
	if err := h.Repos.Events.Delete(ctx, id); err != nil {
		helpers.HandleError(c, err, "Failed to delete event")
		return
	}

	logging.Info(ctx, "event deleted successfully", "event_id", id)
	c.Status(http.StatusNoContent)
}

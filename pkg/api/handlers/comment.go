package handlers

import (
	"net/http"
	"time"

	"github.com/alireza-akbarzadeh/ginflow/pkg/api/helpers"
	"github.com/alireza-akbarzadeh/ginflow/pkg/models"
	"github.com/gin-gonic/gin"
)

// CreateComment handles comment creation
// @Summary      Add a comment to an event
// @Description  Add a comment to an event (requires authentication)
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        id     path      int                 true  "Event ID"
// @Param        comment body      models.Comment  true  "Comment object"
// @Success      201    {object}  models.Comment
// @Failure      400    {object}  helpers.ErrorResponse
// @Failure      401    {object}  helpers.ErrorResponse
// @Failure      404    {object}  helpers.ErrorResponse
// @Failure      500    {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/events/{id}/comments [post]
func (h *Handler) CreateComment(c *gin.Context) {
	eventID, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	// Check if event exists
	event, err := h.Repos.Events.Get(eventID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Database error")
		return
	}
	if event == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "Event not found")
		return
	}

	var comment models.Comment
	if !helpers.BindJSON(c, &comment) {
		return
	}

	// Get authenticated user
	user, ok := helpers.GetAuthenticatedUser(c)
	if !ok {
		return
	}

	comment.EventID = eventID
	comment.UserID = user.ID
	comment.CreatedAt = time.Now()

	createdComment, err := h.Repos.Comments.Insert(&comment)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	c.JSON(http.StatusCreated, createdComment)
}

// GetEventComments retrieves all comments for an event
// @Summary      Get comments for an event
// @Description  Get all comments for a specific event
// @Tags         Comments
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {array}   models.Comment
// @Failure      400  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Router       /api/v1/events/{id}/comments [get]
func (h *Handler) GetEventComments(c *gin.Context) {
	eventID, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	comments, err := h.Repos.Comments.GetByEvent(eventID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}

	c.JSON(http.StatusOK, comments)
}

// DeleteComment removes a comment
// @Summary      Delete a comment
// @Description  Delete a comment by ID (requires authentication)
// @Tags         Comments
// @Param        id         path      int  true  "Event ID"
// @Param        commentId  path      int  true  "Comment ID"
// @Success      204        {object}  nil
// @Failure      400        {object}  helpers.ErrorResponse
// @Failure      401        {object}  helpers.ErrorResponse
// @Failure      403        {object}  helpers.ErrorResponse
// @Failure      404        {object}  helpers.ErrorResponse
// @Failure      500        {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/events/{id}/comments/{commentId} [delete]
func (h *Handler) DeleteComment(c *gin.Context) {
	commentID, err := helpers.ParseIDParam(c, "commentId")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	// Get authenticated user
	user, ok := helpers.GetAuthenticatedUser(c)
	if !ok {
		return
	}

	// Check if comment exists
	comment, err := h.Repos.Comments.Get(commentID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Database error")
		return
	}
	if comment == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "Comment not found")
		return
	}

	// Check ownership (only author can delete)
	if comment.UserID != user.ID {
		helpers.RespondWithError(c, http.StatusForbidden, "You are not allowed to delete this comment")
		return
	}

	if err := h.Repos.Comments.Delete(commentID); err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to delete comment")
		return
	}

	c.Status(http.StatusNoContent)
}

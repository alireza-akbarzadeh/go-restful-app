package handlers

import (
	"net/http"
	"time"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/helpers"
	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
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
	ctx := c.Request.Context()

	eventID, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	logging.Debug(ctx, "creating comment for event", "event_id", eventID)

	// Check if event exists
	_, err = h.Repos.Events.Get(ctx, eventID)
	if helpers.HandleError(c, err, "Failed to retrieve event") {
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

	createdComment, err := h.Repos.Comments.Insert(ctx, &comment)
	if helpers.HandleError(c, err, "Failed to create comment") {
		return
	}

	logging.Info(ctx, "comment created successfully", "comment_id", createdComment.ID, "event_id", eventID, "user_id", user.ID)
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
	ctx := c.Request.Context()

	eventID, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	logging.Debug(ctx, "retrieving comments for event", "event_id", eventID)

	comments, err := h.Repos.Comments.GetByEvent(ctx, eventID)
	if helpers.HandleError(c, err, "Failed to fetch comments") {
		return
	}

	logging.Debug(ctx, "comments retrieved successfully", "event_id", eventID, "count", len(comments))
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
	ctx := c.Request.Context()

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

	logging.Debug(ctx, "deleting comment", "comment_id", commentID, "user_id", user.ID)

	// Check if comment exists
	comment, err := h.Repos.Comments.Get(ctx, commentID)
	if helpers.HandleError(c, err, "Failed to retrieve comment") {
		return
	}
	if comment == nil {
		helpers.RespondWithAppError(c, appErrors.Newf(appErrors.ErrNotFound, "comment with ID %d not found", commentID), "")
		return
	}

	// Check ownership (only author can delete)
	if comment.UserID != user.ID {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrForbidden, "You are not allowed to delete this comment"), "")
		return
	}

	if err := h.Repos.Comments.Delete(ctx, commentID); err != nil {
		helpers.HandleError(c, err, "Failed to delete comment")
		return
	}

	logging.Info(ctx, "comment deleted successfully", "comment_id", commentID, "user_id", user.ID)
	c.Status(http.StatusNoContent)
}

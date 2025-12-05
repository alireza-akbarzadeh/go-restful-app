package repository

import (
	"context"

	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"gorm.io/gorm"
)

// CommentRepository handles comment database operations
type CommentRepository struct {
	DB *gorm.DB
}

// NewCommentRepository creates a new CommentRepository
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

// Insert creates a new comment
func (r *CommentRepository) Insert(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	logging.Debug(ctx, "creating new comment", "event_id", comment.EventID, "user_id", comment.UserID)

	result := r.DB.WithContext(ctx).Create(comment)
	if result.Error != nil {
		logging.Error(ctx, "failed to create comment", result.Error, "event_id", comment.EventID)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to create comment")
	}

	// Preload user info to return complete object
	if err := r.DB.WithContext(ctx).Preload("User").First(comment, comment.ID).Error; err != nil {
		logging.Error(ctx, "failed to preload user data for comment", err, "comment_id", comment.ID)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve complete comment")
	}

	logging.Info(ctx, "comment created successfully", "comment_id", comment.ID, "event_id", comment.EventID)
	return comment, nil
}

// GetByEvent retrieves all comments for a specific event
func (r *CommentRepository) GetByEvent(ctx context.Context, eventID int) ([]*models.Comment, error) {
	var comments []*models.Comment
	result := r.DB.WithContext(ctx).Where("event_id = ?", eventID).Preload("User").Order("created_at desc").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

// Delete removes a comment
func (r *CommentRepository) Delete(ctx context.Context, id int) error {
	result := r.DB.WithContext(ctx).Delete(&models.Comment{}, id)
	return result.Error
}

// Get retrieves a comment by ID
func (r *CommentRepository) Get(ctx context.Context, id int) (*models.Comment, error) {
	var comment models.Comment
	result := r.DB.WithContext(ctx).First(&comment, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &comment, nil
}

package repository

import (
	"context"

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
	result := r.DB.WithContext(ctx).Create(comment)
	if result.Error != nil {
		return nil, result.Error
	}
	// Preload user info to return complete object
	r.DB.WithContext(ctx).Preload("User").First(comment, comment.ID)
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

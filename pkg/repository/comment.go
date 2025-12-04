package repository

import (
	"time"

	"gorm.io/gorm"
)

// Comment represents a comment on an event
type Comment struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"userId" gorm:"not null"`
	EventID   int       `json:"eventId" gorm:"not null"`
	Content   string    `json:"content" binding:"required,min=1" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`

	// Associations (optional, for preloading if needed)
	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Event Event `json:"-" gorm:"foreignKey:EventID"`
}

// CommentRepository handles comment database operations
type CommentRepository struct {
	DB *gorm.DB
}

// NewCommentRepository creates a new CommentRepository
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

// Insert creates a new comment
func (r *CommentRepository) Insert(comment *Comment) (*Comment, error) {
	result := r.DB.Create(comment)
	if result.Error != nil {
		return nil, result.Error
	}
	// Preload user info to return complete object
	r.DB.Preload("User").First(comment, comment.ID)
	return comment, nil
}

// GetByEvent retrieves all comments for a specific event
func (r *CommentRepository) GetByEvent(eventID int) ([]*Comment, error) {
	var comments []*Comment
	result := r.DB.Where("event_id = ?", eventID).Preload("User").Order("created_at desc").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

// Delete removes a comment
func (r *CommentRepository) Delete(id int) error {
	result := r.DB.Delete(&Comment{}, id)
	return result.Error
}

// Get retrieves a comment by ID
func (r *CommentRepository) Get(id int) (*Comment, error) {
	var comment Comment
	result := r.DB.First(&comment, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &comment, nil
}

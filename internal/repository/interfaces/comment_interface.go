package interfaces

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
)

type CommentRepositoryInterface interface {
	Insert(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetByEvent(ctx context.Context, eventID int) ([]*models.Comment, error)
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (*models.Comment, error)
}

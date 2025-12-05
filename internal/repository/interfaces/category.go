package interfaces

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
)

type CategoryRepositoryInterface interface {
	Insert(ctx context.Context, category *models.Category) (*models.Category, error)
	GetAll(ctx context.Context) ([]*models.Category, error)
	Get(ctx context.Context, id int) (*models.Category, error)
	GetBySlug(ctx context.Context, slug string) (*models.Category, error)
}
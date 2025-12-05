package interfaces

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
)

type EventRepositoryInterface interface {
	Insert(ctx context.Context, event *models.Event) (*models.Event, error)
	Get(ctx context.Context, id int) (*models.Event, error)
	GetAll(ctx context.Context) ([]*models.Event, error)
	Update(ctx context.Context, event *models.Event) error
	Delete(ctx context.Context, id int) error
}

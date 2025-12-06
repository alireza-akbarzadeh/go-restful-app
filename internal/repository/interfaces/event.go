package interfaces

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/pagination"
)

type EventRepositoryInterface interface {
	Insert(ctx context.Context, event *models.Event) (*models.Event, error)
	Get(ctx context.Context, id int) (*models.Event, error)
	GetAll(ctx context.Context) ([]*models.Event, error)
	Update(ctx context.Context, event *models.Event) error
	Delete(ctx context.Context, id int) error
	ListWithPagination(ctx context.Context, req *pagination.PaginationRequest) ([]*models.Event, *pagination.PaginationResponse, error)
	ListWithAdvancedPagination(ctx context.Context, req *pagination.AdvancedPaginationRequest) ([]*models.Event, *pagination.AdvancedPaginatedResult, error)
	GetByOwnerID(ctx context.Context, ownerID int) ([]*models.Event, error)
}

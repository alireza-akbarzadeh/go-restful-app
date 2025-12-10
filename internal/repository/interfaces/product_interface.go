package interfaces

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/query"
)

type ProductRepositoryInterface interface {
	Insert(ctx context.Context, product *models.Product) (*models.Product, error)
	GetAll(ctx context.Context, page, limit int, search string, categoryID int) ([]models.Product, int64, error)
	ListWithAdvancedPagination(ctx context.Context, req *query.QueryParams) ([]models.Product, *query.PaginatedList, error)
	Get(ctx context.Context, id int) (*models.Product, error)
	GetBySlug(ctx context.Context, slug string) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
	GetByUser(ctx context.Context, userID int) ([]models.Product, error)
	GetByCategory(ctx context.Context, categoryID int) ([]models.Product, error)
}

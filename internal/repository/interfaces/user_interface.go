package interfaces

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/query"
)

type UserRepositoryInterface interface {
	Insert(ctx context.Context, user *models.User) (*models.User, error)
	Get(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdatePassword(ctx context.Context, userID int, hashedPassword string) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	UpdateLastLogin(ctx context.Context, id int) error
	GetAll(ctx context.Context, params *query.QueryParams) ([]*models.User, *query.PaginatedList, error)
}

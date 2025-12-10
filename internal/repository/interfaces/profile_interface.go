package interfaces

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
)

type ProfileRepositoryInterface interface {
	Insert(ctx context.Context, profile *models.Profile) (*models.Profile, error)
	GetByUserID(ctx context.Context, userID int) (*models.Profile, error)
	GetByUserIDWithUser(ctx context.Context, userID int) (*models.Profile, error)
	Update(ctx context.Context, profile *models.Profile) error
	UpdateByUserID(ctx context.Context, userID int, updates map[string]interface{}) error
	DeleteByUserID(ctx context.Context, id int) error
}

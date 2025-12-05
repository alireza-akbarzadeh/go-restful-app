package interfaces

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
)

type BasketRepositoryInterface interface {
	GetActiveBasket(ctx context.Context, userID int) (*models.Basket, error)
	CreateBasket(ctx context.Context, basket *models.Basket) error
	AddItem(ctx context.Context, basketID int, item *models.BasketItem) error
	UpdateItemQuantity(ctx context.Context, itemID int, quantity int) error
	RemoveItem(ctx context.Context, itemID int) error
	ClearBasket(ctx context.Context, basketID int) error
}

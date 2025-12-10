package repository

import (
	"context"
	"errors"

	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"gorm.io/gorm"
)

type BasketRepository struct {
	DB *gorm.DB
}

func NewBasketRepository(db *gorm.DB) *BasketRepository {
	return &BasketRepository{DB: db}
}

// GetActiveBasket retrieves the active basket for a user
func (r *BasketRepository) GetActiveBasket(ctx context.Context, userID int) (*models.Basket, error) {
	logging.Debug(ctx, "retrieving active basket", "user_id", userID)

	var basket models.Basket
	err := r.DB.WithContext(ctx).Preload("Items.Product").Where("user_id = ? AND status = ?", userID, "active").First(&basket).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logging.Debug(ctx, "no active basket found", "user_id", userID)
			return nil, appErrors.Newf(appErrors.ErrNotFound, "active basket for user ID %d not found", userID)
		}
		logging.Error(ctx, "failed to retrieve active basket", err, "user_id", userID)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve active basket")
	}

	logging.Debug(ctx, "active basket retrieved successfully", "user_id", userID, "basket_id", basket.ID)
	return &basket, nil
}

// CreateBasket creates a new basket
func (r *BasketRepository) CreateBasket(ctx context.Context, basket *models.Basket) error {
	logging.Debug(ctx, "creating new basket", "user_id", basket.UserID)

	if err := r.DB.WithContext(ctx).Create(basket).Error; err != nil {
		logging.Error(ctx, "failed to create basket", err, "user_id", basket.UserID)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to create basket")
	}

	logging.Info(ctx, "basket created successfully", "basket_id", basket.ID, "user_id", basket.UserID)
	return nil
}

// AddItem adds an item to the basket or updates quantity if it exists
func (r *BasketRepository) AddItem(ctx context.Context, basketID int, item *models.BasketItem) error {
	logging.Debug(ctx, "adding item to basket", "basket_id", basketID, "product_id", item.ProductID, "quantity", item.Quantity)

	var existingItem models.BasketItem
	err := r.DB.WithContext(ctx).Where("basket_id = ? AND product_id = ?", basketID, item.ProductID).First(&existingItem).Error

	if err == nil {
		// Item exists, update quantity
		existingItem.Quantity += item.Quantity
		if err := r.DB.WithContext(ctx).Save(&existingItem).Error; err != nil {
			logging.Error(ctx, "failed to update basket item quantity", err, "basket_id", basketID, "product_id", item.ProductID)
			return appErrors.New(appErrors.ErrDatabaseOperation, "failed to update basket item")
		}
		logging.Info(ctx, "basket item quantity updated", "basket_id", basketID, "product_id", item.ProductID, "new_quantity", existingItem.Quantity)
		return nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Item does not exist, create new
		item.BasketID = basketID
		if err := r.DB.WithContext(ctx).Create(item).Error; err != nil {
			logging.Error(ctx, "failed to create basket item", err, "basket_id", basketID, "product_id", item.ProductID)
			return appErrors.New(appErrors.ErrDatabaseOperation, "failed to add item to basket")
		}
		logging.Info(ctx, "item added to basket", "basket_id", basketID, "product_id", item.ProductID, "quantity", item.Quantity)
		return nil
	}

	logging.Error(ctx, "failed to check existing basket item", err, "basket_id", basketID, "product_id", item.ProductID)
	return appErrors.New(appErrors.ErrDatabaseOperation, "failed to process basket item")
}

// UpdateItemQuantity updates the quantity of an item in the basket
func (r *BasketRepository) UpdateItemQuantity(ctx context.Context, itemID int, quantity int) error {
	logging.Debug(ctx, "updating basket item quantity", "item_id", itemID, "quantity", quantity)

	if quantity <= 0 {
		if err := r.DB.WithContext(ctx).Delete(&models.BasketItem{}, itemID).Error; err != nil {
			logging.Error(ctx, "failed to delete basket item", err, "item_id", itemID)
			return appErrors.New(appErrors.ErrDatabaseOperation, "failed to remove basket item")
		}
		logging.Info(ctx, "basket item removed", "item_id", itemID)
		return nil
	}

	result := r.DB.WithContext(ctx).Model(&models.BasketItem{}).Where("id = ?", itemID).Update("quantity", quantity)
	if result.Error != nil {
		logging.Error(ctx, "failed to update basket item quantity", result.Error, "item_id", itemID)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to update basket item quantity")
	}

	if result.RowsAffected == 0 {
		logging.Debug(ctx, "no basket item found to update", "item_id", itemID)
		return appErrors.Newf(appErrors.ErrNotFound, "basket item with ID %d not found", itemID)
	}

	logging.Info(ctx, "basket item quantity updated", "item_id", itemID, "quantity", quantity)
	return nil
}

// RemoveItem removes an item from the basket
func (r *BasketRepository) RemoveItem(ctx context.Context, itemID int) error {
	logging.Debug(ctx, "removing item from basket", "item_id", itemID)

	result := r.DB.WithContext(ctx).Delete(&models.BasketItem{}, itemID)
	if result.Error != nil {
		logging.Error(ctx, "failed to remove basket item", result.Error, "item_id", itemID)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to remove basket item")
	}

	if result.RowsAffected == 0 {
		logging.Debug(ctx, "no basket item found to remove", "item_id", itemID)
		return appErrors.Newf(appErrors.ErrNotFound, "basket item with ID %d not found", itemID)
	}

	logging.Info(ctx, "basket item removed successfully", "item_id", itemID)
	return nil
}

// ClearBasket removes all items from the basket
func (r *BasketRepository) ClearBasket(ctx context.Context, basketID int) error {
	logging.Debug(ctx, "clearing basket", "basket_id", basketID)

	result := r.DB.WithContext(ctx).Where("basket_id = ?", basketID).Delete(&models.BasketItem{})
	if result.Error != nil {
		logging.Error(ctx, "failed to clear basket", result.Error, "basket_id", basketID)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to clear basket")
	}

	logging.Info(ctx, "basket cleared successfully", "basket_id", basketID, "items_removed", result.RowsAffected)
	return nil
}

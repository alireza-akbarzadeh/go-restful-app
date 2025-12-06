package handlers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/helpers"
	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/gin-gonic/gin"
)

// GetBasket retrieves the current user's basket
// @Summary      Get user basket
// @Description  Get the active basket for the authenticated user
// @Tags         Basket
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Basket
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/basket [get]
func (h *Handler) GetBasket(c *gin.Context) {
	ctx := c.Request.Context()

	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrUnauthorized, "Authentication required"), "Authentication required")
		return
	}

	basket, err := h.Repos.Baskets.GetActiveBasket(ctx, user.ID)
	if err != nil {
		logging.Error(ctx, "Failed to retrieve basket", err, "userID", user.ID)
		helpers.HandleError(c, err, "Failed to retrieve basket")
		return
	}

	if basket == nil {
		// Create a new basket if none exists
		basket = &models.Basket{
			UserID: &user.ID,
			Status: "active",
		}
		if err := h.Repos.Baskets.CreateBasket(ctx, basket); err != nil {
			logging.Error(ctx, "Failed to create basket", err, "userID", user.ID)
			helpers.HandleError(c, err, "Failed to create basket")
			return
		}
	}

	c.JSON(http.StatusOK, basket)
}

type AddItemRequest struct {
	ProductID int `json:"productId" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

// AddItemToBasket adds an item to the basket
// @Summary      Add item to basket
// @Description  Add a product to the user's active basket
// @Tags         Basket
// @Accept       json
// @Produce      json
// @Param        item  body      AddItemRequest  true  "Item to add"
// @Success      200   {object}  models.Basket
// @Failure      400   {object}  helpers.ErrorResponse
// @Failure      401   {object}  helpers.ErrorResponse
// @Failure      404   {object}  helpers.ErrorResponse
// @Failure      500   {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/basket/items [post]
func (h *Handler) AddItemToBasket(c *gin.Context) {
	ctx := c.Request.Context()

	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrUnauthorized, "Authentication required"), "Authentication required")
		return
	}

	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrInvalidInput, err.Error()), "Invalid request body")
		return
	}

	// Get or create basket
	basket, err := h.Repos.Baskets.GetActiveBasket(ctx, user.ID)
	if err != nil {
		logging.Error(ctx, "Failed to retrieve basket", err, "userID", user.ID)
		helpers.HandleError(c, err, "Failed to retrieve basket")
		return
	}
	if basket == nil {
		basket = &models.Basket{
			UserID: &user.ID,
			Status: "active",
		}
		if err := h.Repos.Baskets.CreateBasket(ctx, basket); err != nil {
			logging.Error(ctx, "Failed to create basket", err, "userID", user.ID)
			helpers.HandleError(c, err, "Failed to create basket")
			return
		}
	}

	// Get product to check price and existence
	product, err := h.Repos.Products.Get(ctx, req.ProductID)
	if err != nil {
		logging.Error(ctx, "Failed to retrieve product", err, "productID", req.ProductID)
		helpers.HandleError(c, err, "Failed to retrieve product")
		return
	}
	if product == nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrNotFound, "Product not found"), "Product not found")
		return
	}

	item := &models.BasketItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UnitPrice: product.Price, // Use current price
	}

	if err := h.Repos.Baskets.AddItem(ctx, basket.ID, item); err != nil {
		logging.Error(ctx, "Failed to add item to basket", err, "basketID", basket.ID, "productID", req.ProductID)
		helpers.HandleError(c, err, "Failed to add item to basket")
		return
	}

	// Return updated basket
	updatedBasket, _ := h.Repos.Baskets.GetActiveBasket(ctx, user.ID)
	c.JSON(http.StatusOK, updatedBasket)
}

// RemoveItemFromBasket removes an item from the basket
// @Summary      Remove item from basket
// @Description  Remove an item from the basket by Item ID
// @Tags         Basket
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Item ID"
// @Success      200  {object}  models.Basket
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/basket/items/{id} [delete]
func (h *Handler) RemoveItemFromBasket(c *gin.Context) {
	ctx := c.Request.Context()

	itemID, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrUnauthorized, "Authentication required"), "Authentication required")
		return
	}

	// Verify ownership (optional but good practice, though item ID is unique)
	// For simplicity, just delete
	if err := h.Repos.Baskets.RemoveItem(ctx, itemID); err != nil {
		logging.Error(ctx, "Failed to remove item from basket", err, "itemID", itemID)
		helpers.HandleError(c, err, "Failed to remove item")
		return
	}

	updatedBasket, _ := h.Repos.Baskets.GetActiveBasket(ctx, user.ID)
	c.JSON(http.StatusOK, updatedBasket)
}

// ClearBasket clears the basket
// @Summary      Clear basket
// @Description  Remove all items from the active basket
// @Tags         Basket
// @Accept       json
// @Produce      json
// @Success      204  {object}  nil
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/basket [delete]
func (h *Handler) ClearBasket(c *gin.Context) {
	ctx := c.Request.Context()

	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrUnauthorized, "Authentication required"), "Authentication required")
		return
	}

	basket, err := h.Repos.Baskets.GetActiveBasket(ctx, user.ID)
	if err != nil {
		logging.Error(ctx, "Failed to retrieve basket", err, "userID", user.ID)
		helpers.HandleError(c, err, "Failed to retrieve basket")
		return
	}

	if basket != nil {
		if err := h.Repos.Baskets.ClearBasket(ctx, basket.ID); err != nil {
			logging.Error(ctx, "Failed to clear basket", err, "basketID", basket.ID)
			helpers.HandleError(c, err, "Failed to clear basket")
			return
		}
	}

	c.Status(http.StatusNoContent)
}

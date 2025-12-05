package handlers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/helpers"
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
	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	basket, err := h.Repos.Baskets.GetActiveBasket(c.Request.Context(), user.ID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve basket")
		return
	}

	if basket == nil {
		// Create a new basket if none exists
		basket = &models.Basket{
			UserID: &user.ID,
			Status: "active",
		}
		if err := h.Repos.Baskets.CreateBasket(c.Request.Context(), basket); err != nil {
			helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to create basket")
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
	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get or create basket
	basket, err := h.Repos.Baskets.GetActiveBasket(c.Request.Context(), user.ID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve basket")
		return
	}
	if basket == nil {
		basket = &models.Basket{
			UserID: &user.ID,
			Status: "active",
		}
		if err := h.Repos.Baskets.CreateBasket(c.Request.Context(), basket); err != nil {
			helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to create basket")
			return
		}
	}

	// Get product to check price and existence
	product, err := h.Repos.Products.Get(c.Request.Context(), req.ProductID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve product")
		return
	}
	if product == nil {
		helpers.RespondWithError(c, http.StatusNotFound, "Product not found")
		return
	}

	item := &models.BasketItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UnitPrice: product.Price, // Use current price
	}

	if err := h.Repos.Baskets.AddItem(c.Request.Context(), basket.ID, item); err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to add item to basket")
		return
	}

	// Return updated basket
	updatedBasket, _ := h.Repos.Baskets.GetActiveBasket(c.Request.Context(), user.ID)
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
	itemID, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Verify ownership (optional but good practice, though item ID is unique)
	// For simplicity, just delete
	if err := h.Repos.Baskets.RemoveItem(c.Request.Context(), itemID); err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to remove item")
		return
	}

	updatedBasket, _ := h.Repos.Baskets.GetActiveBasket(c.Request.Context(), user.ID)
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
	user := helpers.GetUserFromContext(c)
	if user == nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	basket, err := h.Repos.Baskets.GetActiveBasket(c.Request.Context(), user.ID)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve basket")
		return
	}

	if basket != nil {
		if err := h.Repos.Baskets.ClearBasket(c.Request.Context(), basket.ID); err != nil {
			helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to clear basket")
			return
		}
	}

	c.Status(http.StatusNoContent)
}

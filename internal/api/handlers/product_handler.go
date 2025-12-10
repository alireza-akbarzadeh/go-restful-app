package handlers

import (
	"net/http"
	"strconv"

	"github.com/alireza-akbarzadeh/ginflow/internal/api/helpers"
	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/query"
	"github.com/alireza-akbarzadeh/ginflow/internal/utils"
	"github.com/gin-gonic/gin"
)

// CreateProduct handles product creation
// @Summary      Create a new product
// @Description  Create a new product (requires authentication)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "Product object"
// @Success      201      {object}  models.Product
// @Failure      400      {object}  helpers.ErrorResponse
// @Failure      401      {object}  helpers.ErrorResponse
// @Failure      500      {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/products [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	var product models.Product
	if !helpers.BindJSON(c, &product) {
		return
	}

	// Get authenticated user
	user, ok := helpers.GetAuthenticatedUser(c)
	if !ok {
		return
	}
	product.UserID = user.ID

	logging.Debug(ctx, "creating new product", "name", product.Name, "user_id", user.ID)

	// Generate slug if not provided
	if product.Slug == "" {
		product.Slug = utils.GenerateSlug(product.Name)
	}

	createdProduct, err := h.Repos.Products.Insert(ctx, &product)
	if helpers.HandleError(c, err, "Failed to create product") {
		return
	}

	logging.Info(ctx, "product created successfully", "product_id", createdProduct.ID, "name", createdProduct.Name)
	c.JSON(http.StatusCreated, createdProduct)
}

// GetAllProducts retrieves all products with advanced pagination
// @Summary      Get all products
// @Description  Get a list of all products with advanced pagination, filtering, sorting, and search
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        page        query     int     false  "Page number (default: 1)"
// @Param        page_size   query     int     false  "Page size (default: 20, max: 100)"
// @Param        type        query     string  false  "Pagination type: 'offset' or 'cursor' (default: offset)"
// @Param        cursor      query     string  false  "Cursor for cursor-based pagination"
// @Param        sort        query     string  false  "Sort fields (e.g., '-created_at,name:asc,price:desc')"
// @Param        search      query     string  false  "Search term for name, slug, description"
// @Param        name[eq]    query     string  false  "Filter by exact name"
// @Param        name[like]  query     string  false  "Filter by name (partial match)"
// @Param        price[gte]  query     number  false  "Filter by minimum price"
// @Param        price[lte]  query     number  false  "Filter by maximum price"
// @Param        user_id[eq] query     int     false  "Filter by user ID"
// @Success      200         {object}  query.PaginatedList{data=[]models.Product}
// @Failure      500         {object}  helpers.ErrorResponse
// @Router       /api/v1/products [get]
func (h *Handler) GetAllProducts(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse advanced pagination parameters from context
	req := query.ParseFromContext(c)

	logging.Debug(ctx, "retrieving all products with advanced pagination",
		"page", req.Page,
		"page_size", req.PageSize,
		"type", req.Type,
		"search", req.Search,
	)

	products, result, err := h.Repos.Products.ListWithAdvancedPagination(ctx, req)
	if helpers.HandleError(c, err, "Failed to retrieve products") {
		return
	}

	logging.Debug(ctx, "products retrieved successfully",
		"count", len(products),
		"page", req.Page,
	)

	c.JSON(http.StatusOK, result)
}

// GetProduct retrieves a product by ID or Slug
// @Summary      Get a product
// @Description  Get a product by ID or Slug
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID or Slug"
// @Success      200  {object}  models.Product
// @Failure      404  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Router       /api/v1/products/{id} [get]
func (h *Handler) GetProduct(c *gin.Context) {
	ctx := c.Request.Context()
	idOrSlug := c.Param("id")

	var product *models.Product
	var err error

	// Try to parse as ID first
	if id, parseErr := strconv.Atoi(idOrSlug); parseErr == nil {
		product, err = h.Repos.Products.Get(ctx, id)
	} else {
		product, err = h.Repos.Products.GetBySlug(ctx, idOrSlug)
	}

	if helpers.HandleError(c, err, "Failed to retrieve product") {
		return
	}
	if product == nil {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrNotFound, "Product not found"), "")
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct updates a product
// @Summary      Update a product
// @Description  Update a product by ID (Owner only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id       path      int             true  "Product ID"
// @Param        product  body      models.Product  true  "Product object"
// @Success      200      {object}  models.Product
// @Failure      400      {object}  helpers.ErrorResponse
// @Failure      401      {object}  helpers.ErrorResponse
// @Failure      403      {object}  helpers.ErrorResponse
// @Failure      404      {object}  helpers.ErrorResponse
// @Failure      500      {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/products/{id} [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	user, ok := helpers.GetAuthenticatedUser(c)
	if !ok {
		return
	}

	logging.Debug(ctx, "updating product", "product_id", id, "user_id", user.ID)

	existingProduct, err := h.Repos.Products.Get(ctx, id)
	if helpers.HandleError(c, err, "Failed to retrieve product") {
		return
	}
	if existingProduct == nil {
		helpers.RespondWithAppError(c, appErrors.Newf(appErrors.ErrNotFound, "product with ID %d not found", id), "")
		return
	}
	if existingProduct.UserID != user.ID {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrForbidden, "You do not have permission to update this product"), "")
		return
	}

	var updateData models.Product
	if err := c.ShouldBindJSON(&updateData); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update fields
	existingProduct.Name = updateData.Name
	existingProduct.Description = updateData.Description
	existingProduct.Price = updateData.Price
	existingProduct.Stock = updateData.Stock
	existingProduct.SKU = updateData.SKU
	existingProduct.Status = updateData.Status
	existingProduct.Image = updateData.Image
	existingProduct.Images = updateData.Images
	existingProduct.Tags = updateData.Tags
	existingProduct.MetaTitle = updateData.MetaTitle
	existingProduct.MetaDescription = updateData.MetaDescription
	existingProduct.Discount = updateData.Discount
	existingProduct.FinalPrice = updateData.FinalPrice
	existingProduct.Brand = updateData.Brand
	existingProduct.Weight = updateData.Weight
	existingProduct.Dimensions = updateData.Dimensions

	if err := h.Repos.Products.Update(ctx, existingProduct); err != nil {
		helpers.HandleError(c, err, "Failed to update product")
		return
	}

	logging.Info(ctx, "product updated successfully", "product_id", id, "name", existingProduct.Name)
	c.JSON(http.StatusOK, existingProduct)
}

// DeleteProduct deletes a product
// @Summary      Delete a product
// @Description  Delete a product by ID (Owner only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      204  {object}  nil
// @Failure      401  {object}  helpers.ErrorResponse
// @Failure      403  {object}  helpers.ErrorResponse
// @Failure      404  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/products/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	user, ok := helpers.GetAuthenticatedUser(c)
	if !ok {
		return
	}

	logging.Debug(ctx, "deleting product", "product_id", id, "user_id", user.ID)

	existingProduct, err := h.Repos.Products.Get(ctx, id)
	if helpers.HandleError(c, err, "Failed to retrieve product") {
		return
	}
	if existingProduct == nil {
		helpers.RespondWithAppError(c, appErrors.Newf(appErrors.ErrNotFound, "product with ID %d not found", id), "")
		return
	}

	// Check ownership
	if existingProduct.UserID != user.ID {
		helpers.RespondWithAppError(c, appErrors.New(appErrors.ErrForbidden, "You are not allowed to delete this product"), "")
		return
	}

	if err := h.Repos.Products.Delete(ctx, id); err != nil {
		helpers.HandleError(c, err, "Failed to delete product")
		return
	}

	logging.Info(ctx, "product deleted successfully", "product_id", id)
	c.Status(http.StatusNoContent)
}

// GetProductBySlug retrieves a product by Slug
// @Summary      Get a product by Slug
// @Description  Get a product by Slug
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        slug   path      string  true  "Product Slug"
// @Success      200  {object}  models.Product
// @Failure      404  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Router       /api/v1/products/slug/{slug} [get]
func (h *Handler) GetProductBySlug(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	logging.Debug(ctx, "retrieving product by slug", "slug", slug)

	product, err := h.Repos.Products.GetBySlug(ctx, slug)
	if helpers.HandleError(c, err, "Failed to retrieve product") {
		return
	}
	if product == nil {
		helpers.RespondWithAppError(c, appErrors.Newf(appErrors.ErrNotFound, "product with slug '%s' not found", slug), "")
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProductsByCategory retrieves products by Category ID
// @Summary      Get products by Category
// @Description  Get products by Category ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  []models.Product
// @Failure      500  {object}  helpers.ErrorResponse
// @Router       /api/v1/products/category/{id} [get]
func (h *Handler) GetProductsByCategory(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	logging.Debug(ctx, "retrieving products by category", "category_id", id)

	products, err := h.Repos.Products.GetByCategory(ctx, id)
	if helpers.HandleError(c, err, "Failed to retrieve products") {
		return
	}

	logging.Debug(ctx, "products retrieved by category", "category_id", id, "count", len(products))
	c.JSON(http.StatusOK, products)
}

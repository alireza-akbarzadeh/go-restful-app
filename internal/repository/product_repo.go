package repository

import (
	"context"
	"errors"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/query"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) Insert(ctx context.Context, product *models.Product) (*models.Product, error) {
	result := r.DB.WithContext(ctx).Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

// GetAll retrieves all products with optional pagination and filters
func (r *ProductRepository) GetAll(ctx context.Context, page, limit int, search string, categoryID int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := (page - 1) * limit
	query := r.DB.WithContext(ctx).Model(&models.Product{})

	// Apply filters
	if search != "" {
		query = query.Where("name ILIKE ? OR slug ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if categoryID > 0 {
		query = query.Joins("JOIN product_categories ON products.id = product_categories.product_id").
			Where("product_categories.category_id = ?", categoryID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch products with relationships
	result := query.Preload("User").Preload("Categories").
		Offset(offset).Limit(limit).
		Order("created_at desc").
		Find(&products)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return products, total, nil
}

// ListWithAdvancedPagination retrieves products with advanced pagination, filtering, sorting, and search
func (r *ProductRepository) ListWithAdvancedPagination(ctx context.Context, req *query.QueryParams) ([]models.Product, *query.PaginatedList, error) {
	var products []models.Product
	var total int64

	// Build pagination query
	builder := query.NewQueryBuilder(r.DB.WithContext(ctx).Model(&models.Product{})).
		WithRequest(req).
		AllowFilters("name", "slug", "user_id", "price", "status", "created_at").
		AllowSorts("name", "price", "created_at", "updated_at").
		SearchColumns("name", "slug", "description").
		DefaultSort("created_at", query.SortDesc)

	// Get count if needed
	if req.IncludeTotal {
		countQuery := r.DB.WithContext(ctx).Model(&models.Product{})
		for _, filter := range req.Filters {
			countQuery = query.FilterBy(filter)(countQuery)
		}
		if req.Search != "" {
			countQuery = query.Search(req.Search, "name", "slug", "description")(countQuery)
		}
		countQuery.Count(&total)
	}

	// Execute main query
	dbQuery := builder.Build()
	if err := dbQuery.Preload("User").Preload("Categories").Find(&products).Error; err != nil {
		return nil, nil, err
	}

	// Get first and last IDs for cursor pagination
	var firstID, lastID int
	if len(products) > 0 {
		firstID = products[0].ID
		lastID = products[len(products)-1].ID
	}

	// Build response
	result := query.BuildResponse(products, req, total, len(products), firstID, lastID)

	return products, result, nil
}

// Get retrieves a product by ID
func (r *ProductRepository) Get(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	result := r.DB.WithContext(ctx).Preload("User").Preload("Categories").First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &product, nil
}

// GetBySlug retrieves a product by its slug
func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var product models.Product
	result := r.DB.WithContext(ctx).Preload("User").Preload("Categories").Where("slug = ?", slug).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &product, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	result := r.DB.WithContext(ctx).Save(product)
	return result.Error
}

// Delete removes a product by ID
func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	result := r.DB.WithContext(ctx).Delete(&models.Product{}, id)
	return result.Error
}

// GetByUser retrieves all products created by a specific user
func (r *ProductRepository) GetByUser(ctx context.Context, userID int) ([]models.Product, error) {
	var products []models.Product
	result := r.DB.WithContext(ctx).Where("user_id = ?", userID).Preload("Categories").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

// GetByCategory retrieves all products in a specific category
func (r *ProductRepository) GetByCategory(ctx context.Context, categoryID int) ([]models.Product, error) {
	var products []models.Product
	result := r.DB.WithContext(ctx).Joins("JOIN product_categories ON products.id = product_categories.product_id").
		Where("product_categories.category_id = ?", categoryID).
		Preload("User").Preload("Categories").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

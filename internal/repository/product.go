package repository

import (
	"context"
	"errors"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
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

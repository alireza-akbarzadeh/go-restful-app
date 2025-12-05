package repository

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository handles category database operations
type CategoryRepository struct {
	DB *gorm.DB
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

// Insert creates a new category
func (r *CategoryRepository) Insert(ctx context.Context, category *models.Category) (*models.Category, error) {
	result := r.DB.WithContext(ctx).Create(category)
	if result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

// GetAll retrieves all categories
func (r *CategoryRepository) GetAll(ctx context.Context) ([]*models.Category, error) {
	var categories []*models.Category
	result := r.DB.WithContext(ctx).Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

// Get retrieves a category by ID
func (r *CategoryRepository) Get(ctx context.Context, id int) (*models.Category, error) {
	var category models.Category
	result := r.DB.WithContext(ctx).Preload("Children").Preload("Parent").First(&category, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &category, nil
}

// GetBySlug retrieves a category by slug
func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var category models.Category
	result := r.DB.WithContext(ctx).Preload("Children").Preload("Parent").Where("slug = ?", slug).First(&category)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &category, nil
}

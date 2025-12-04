package repository

import (
	"gorm.io/gorm"
)

// Category represents an event category
type Category struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" binding:"required,min=3" gorm:"uniqueIndex;not null"`
	Description string `json:"description"`
}

// CategoryRepository handles category database operations
type CategoryRepository struct {
	DB *gorm.DB
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

// Insert creates a new category
func (r *CategoryRepository) Insert(category *Category) (*Category, error) {
	result := r.DB.Create(category)
	if result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

// GetAll retrieves all categories
func (r *CategoryRepository) GetAll() ([]*Category, error) {
	var categories []*Category
	result := r.DB.Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

// Get retrieves a category by ID
func (r *CategoryRepository) Get(id int) (*Category, error) {
	var category Category
	result := r.DB.First(&category, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &category, nil
}

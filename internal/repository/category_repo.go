package repository

import (
	"context"
	"errors"

	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/query"
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
	logging.Debug(ctx, "creating new category", "name", category.Name)

	result := r.DB.WithContext(ctx).Create(category)
	if result.Error != nil {
		logging.Error(ctx, "failed to create category", result.Error, "name", category.Name)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to create category")
	}

	logging.Info(ctx, "category created successfully", "category_id", category.ID, "name", category.Name)
	return category, nil
}

// GetAll retrieves all categories
func (r *CategoryRepository) GetAll(ctx context.Context) ([]*models.Category, error) {
	logging.Debug(ctx, "retrieving all categories")

	var categories []*models.Category
	result := r.DB.WithContext(ctx).Find(&categories)
	if result.Error != nil {
		logging.Error(ctx, "failed to retrieve categories", result.Error)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve categories")
	}

	logging.Info(ctx, "categories retrieved successfully", "count", len(categories))
	return categories, nil
}

// Get retrieves a category by ID
func (r *CategoryRepository) Get(ctx context.Context, id int) (*models.Category, error) {
	logging.Debug(ctx, "retrieving category by ID", "category_id", id)

	var category models.Category
	result := r.DB.WithContext(ctx).Preload("Children").Preload("Parent").First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logging.Debug(ctx, "category not found", "category_id", id)
			return nil, appErrors.Newf(appErrors.ErrNotFound, "category with ID %d not found", id)
		}
		logging.Error(ctx, "failed to retrieve category", result.Error, "category_id", id)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve category")
	}

	logging.Debug(ctx, "category retrieved successfully", "category_id", id, "name", category.Name)
	return &category, nil
}

// GetBySlug retrieves a category by slug
func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	logging.Debug(ctx, "retrieving category by slug", "slug", slug)

	var category models.Category
	result := r.DB.WithContext(ctx).Preload("Children").Preload("Parent").Where("slug = ?", slug).First(&category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logging.Debug(ctx, "category not found by slug", "slug", slug)
			return nil, appErrors.Newf(appErrors.ErrNotFound, "category with slug '%s' not found", slug)
		}
		logging.Error(ctx, "failed to retrieve category by slug", result.Error, "slug", slug)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve category by slug")
	}

	logging.Debug(ctx, "category retrieved by slug successfully", "slug", slug, "name", category.Name)
	return &category, nil
}

// ListWithPagination retrieves categories with pagination
func (r *CategoryRepository) ListWithPagination(ctx context.Context, req *query.PaginationRequest) ([]*models.Category, *query.PaginationResponse, error) {
	logging.Debug(ctx, "retrieving categories with pagination", "page", req.Page, "page_size", req.PageSize)

	var categories []*models.Category
	var total int64

	// Count total records
	if err := r.DB.WithContext(ctx).Model(&models.Category{}).Count(&total).Error; err != nil {
		logging.Error(ctx, "failed to count categories", err)
		return nil, nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to count categories")
	}

	// Get paginated records
	if err := r.DB.WithContext(ctx).
		Preload("Parent").
		Offset(req.Offset()).
		Limit(req.PageSize).
		Find(&categories).Error; err != nil {
		logging.Error(ctx, "failed to retrieve categories", err)
		return nil, nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve categories")
	}

	// Calculate pagination response
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	paginationResp := &query.PaginationResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}

	logging.Info(ctx, "categories retrieved successfully", "count", len(categories), "total", total, "page", req.Page)
	return categories, paginationResp, nil
}

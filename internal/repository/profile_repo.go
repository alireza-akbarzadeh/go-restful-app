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

// ProfileRepository handles profile database operations
type ProfileRepository struct {
	DB *gorm.DB
}

// NewProfileRepository creates a new ProfileRepository
func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{DB: db}
}

// Insert creates a new profile in the database
func (r *ProfileRepository) Insert(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	logging.Debug(ctx, "creating new profile", "user_id", profile.UserID)

	result := r.DB.WithContext(ctx).Create(profile)
	if result.Error != nil {
		logging.Error(ctx, "failed to create profile", result.Error, "user_id", profile.UserID)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to create profile")
	}

	logging.Info(ctx, "profile created successfully", "profile_id", profile.ID, "user_id", profile.UserID)
	return profile, nil
}

// Get retrieves a profile by ID
func (r *ProfileRepository) Get(ctx context.Context, id int) (*models.Profile, error) {
	logging.Debug(ctx, "retrieving profile by ID", "profile_id", id)

	var profile models.Profile
	result := r.DB.WithContext(ctx).First(&profile, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logging.Debug(ctx, "profile not found", "profile_id", id)
			return nil, appErrors.Newf(appErrors.ErrNotFound, "profile with ID %d not found", id)
		}
		logging.Error(ctx, "failed to retrieve profile", result.Error, "profile_id", id)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve profile")
	}

	logging.Debug(ctx, "profile retrieved successfully", "profile_id", id)
	return &profile, nil
}

// GetByUserID retrieves a profile by user ID
func (r *ProfileRepository) GetByUserID(ctx context.Context, userID int) (*models.Profile, error) {
	logging.Debug(ctx, "retrieving profile by user ID", "user_id", userID)

	var profile models.Profile
	result := r.DB.WithContext(ctx).Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logging.Debug(ctx, "profile not found for user", "user_id", userID)
			return nil, appErrors.Newf(appErrors.ErrNotFound, "profile for user ID %d not found", userID)
		}
		logging.Error(ctx, "failed to retrieve profile by user ID", result.Error, "user_id", userID)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve profile")
	}

	logging.Debug(ctx, "profile retrieved successfully", "user_id", userID, "profile_id", profile.ID)
	return &profile, nil
}

// GetByUserIDWithUser retrieves a profile with user data preloaded
func (r *ProfileRepository) GetByUserIDWithUser(ctx context.Context, userID int) (*models.Profile, error) {
	logging.Debug(ctx, "retrieving profile with user data", "user_id", userID)

	var profile models.Profile
	result := r.DB.WithContext(ctx).Preload("User").Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logging.Debug(ctx, "profile with user data not found", "user_id", userID)
			return nil, appErrors.Newf(appErrors.ErrNotFound, "profile with user data for user ID %d not found", userID)
		}
		logging.Error(ctx, "failed to retrieve profile with user data", result.Error, "user_id", userID)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve profile with user data")
	}

	logging.Debug(ctx, "profile with user data retrieved successfully", "user_id", userID, "profile_id", profile.ID)
	return &profile, nil
}

// Update updates an existing profile
func (r *ProfileRepository) Update(ctx context.Context, profile *models.Profile) error {
	logging.Debug(ctx, "updating profile", "profile_id", profile.ID, "user_id", profile.UserID)

	result := r.DB.WithContext(ctx).Save(profile)
	if result.Error != nil {
		logging.Error(ctx, "failed to update profile", result.Error, "profile_id", profile.ID)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to update profile")
	}

	logging.Info(ctx, "profile updated successfully", "profile_id", profile.ID, "user_id", profile.UserID)
	return nil
}

// UpdateByUserID updates a profile by user ID
func (r *ProfileRepository) UpdateByUserID(ctx context.Context, userID int, updates map[string]interface{}) error {
	logging.Debug(ctx, "updating profile by user ID", "user_id", userID, "updates", updates)

	result := r.DB.WithContext(ctx).Model(&models.Profile{}).Where("user_id = ?", userID).Updates(updates)
	if result.Error != nil {
		logging.Error(ctx, "failed to update profile by user ID", result.Error, "user_id", userID)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to update profile")
	}

	if result.RowsAffected == 0 {
		logging.Debug(ctx, "no profile found to update", "user_id", userID)
		return appErrors.Newf(appErrors.ErrNotFound, "profile for user ID %d not found", userID)
	}

	logging.Info(ctx, "profile updated successfully by user ID", "user_id", userID, "rows_affected", result.RowsAffected)
	return nil
}

// DeleteByUserID deletes a profile by user ID
func (r *ProfileRepository) DeleteByUserID(ctx context.Context, id int) error {
	logging.Debug(ctx, "deleting profile by user ID", "user_id", id)

	result := r.DB.WithContext(ctx).Delete(&models.Profile{}, id)
	if result.Error != nil {
		logging.Error(ctx, "failed to delete profile", result.Error, "user_id", id)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to delete profile")
	}

	if result.RowsAffected == 0 {
		logging.Debug(ctx, "no profile found to delete", "user_id", id)
		return appErrors.Newf(appErrors.ErrNotFound, "profile for user ID %d not found", id)
	}

	logging.Info(ctx, "profile deleted successfully", "user_id", id)
	return nil
}

// ListWithPagination retrieves profiles with pagination
func (r *ProfileRepository) ListWithPagination(ctx context.Context, req *query.PaginationRequest) ([]*models.Profile, *query.PaginationResponse, error) {
	logging.Debug(ctx, "retrieving profiles with pagination", "page", req.Page, "page_size", req.PageSize)

	var profiles []*models.Profile
	var total int64

	// Count total records
	if err := r.DB.WithContext(ctx).Model(&models.Profile{}).Count(&total).Error; err != nil {
		logging.Error(ctx, "failed to count profiles", err)
		return nil, nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to count profiles")
	}

	// Get paginated records
	if err := r.DB.WithContext(ctx).
		Offset(req.Offset()).
		Limit(req.PageSize).
		Find(&profiles).Error; err != nil {
		logging.Error(ctx, "failed to retrieve profiles", err)
		return nil, nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve profiles")
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

	logging.Info(ctx, "profiles retrieved successfully", "count", len(profiles), "total", total, "page", req.Page)
	return profiles, paginationResp, nil
}

// SearchWithPagination searches profiles with pagination
func (r *ProfileRepository) SearchWithPagination(ctx context.Context, searchTerm string, req *query.PaginationRequest) ([]*models.Profile, *query.PaginationResponse, error) {
	logging.Debug(ctx, "searching profiles with pagination", "search_term", searchTerm, "page", req.Page, "page_size", req.PageSize)

	var profiles []*models.Profile
	var total int64

	dbQuery := r.DB.WithContext(ctx).Model(&models.Profile{}).
		Where("first_name ILIKE ? OR last_name ILIKE ? OR bio ILIKE ?",
			"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		logging.Error(ctx, "failed to count search results", err, "search_term", searchTerm)
		return nil, nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to count search results")
	}

	// Get paginated records
	if err := dbQuery.Offset(req.Offset()).
		Limit(req.PageSize).
		Find(&profiles).Error; err != nil {
		logging.Error(ctx, "failed to search profiles", err, "search_term", searchTerm)
		return nil, nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to search profiles")
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

	logging.Info(ctx, "profile search completed", "results_count", len(profiles), "total", total, "search_term", searchTerm)
	return profiles, paginationResp, nil
}

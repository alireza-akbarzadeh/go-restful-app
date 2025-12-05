package repository

import (
	"context"
	"errors"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
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
	result := r.DB.WithContext(ctx).Create(profile)
	if result.Error != nil {
		return nil, result.Error
	}
	return profile, nil
}

// Get retrieves a profile by ID
func (r *ProfileRepository) Get(ctx context.Context, id int) (*models.Profile, error) {
	var profile models.Profile
	result := r.DB.WithContext(ctx).First(&profile, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &profile, nil
}

// GetByUserID retrieves a profile by user ID

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID int) (*models.Profile, error) {
	var profile models.Profile
	result := r.DB.WithContext(ctx).Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &profile, nil
}

// GetByUserIDWithUser retrieves a profile with user data preloaded

func (r *ProfileRepository) GetByUserIDWithUser(ctx context.Context, userID int) (*models.Profile, error) {
	var profile models.Profile
	result := r.DB.WithContext(ctx).Preload("User").Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &profile, nil
}

// Update updates an existing profile

func (r *ProfileRepository) Update(ctx context.Context, profile *models.Profile) error {
	result := r.DB.WithContext(ctx).Save(profile)
	return result.Error
}

// UpdateByUserID updates a profile by user ID
func (r *ProfileRepository) UpdateByUserID(ctx context.Context, userID int, updates map[string]interface{}) error {
	result := r.DB.WithContext(ctx).Model(&models.Profile{}).Where("user_id = ?", userID).Updates(updates)
	return result.Error
}

// DeleteByUserID deletes a profile by user ID
func (r *ProfileRepository) DeleteByUserID(ctx context.Context, id int) error {
	result := r.DB.WithContext(ctx).Delete(&models.Profile{}, id)
	return result.Error
}

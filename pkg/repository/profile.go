package repository

import (
	"errors"

	"github.com/alireza-akbarzadeh/ginflow/pkg/models"
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
func (r *ProfileRepository) Insert(profile *models.Profile) (*models.Profile, error) {
	result := r.DB.Create(profile)
	if result.Error != nil {
		return nil, result.Error
	}
	return profile, nil
}

// Get retrieves a profile by ID
func (r *ProfileRepository) Get(id uint) (*models.Profile, error) {
	var profile models.Profile
	result := r.DB.First(&profile, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &profile, nil
}

// GetByUserID retrieves a profile by user ID

func (r *ProfileRepository) GetByUserID(userID int) (*models.Profile, error) {
	var profile models.Profile
	result := r.DB.Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &profile, nil
}

// GetByUserIDWithUser retrieves a profile with user data preloaded

func (r *ProfileRepository) GetByUserIDWithUser(userID int) (*models.Profile, error) {
	var profile models.Profile
	result := r.DB.Preload("User").Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &profile, nil
}

// Update updates an existing profile

func (r *ProfileRepository) Update(profile *models.Profile) error {
	result := r.DB.Save(profile)
	return result.Error
}

// UpdateByUserID updates a profile by user ID
func (r *ProfileRepository) UpdateByUserID(userID int, updates map[string]interface{}) error {
	result := r.DB.Model(&models.Profile{}).Where("user_id = ?", userID).Updates(updates)
	return result.Error
}

// DeleteByUserID deletes a profile by user ID
func (r *ProfileRepository) DeleteByUserID(id int) error {
	result := r.DB.Delete(&models.Profile{}, id)
	return result.Error
}

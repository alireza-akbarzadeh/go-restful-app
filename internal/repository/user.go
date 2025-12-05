package repository

import (
	"context"
	"errors"
	"time"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles user database operations
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Insert creates a new user in the database
func (r *UserRepository) Insert(ctx context.Context, user *models.User) (*models.User, error) {
	result := r.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// Get retrieves a user by ID
func (r *UserRepository) Get(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	result := r.DB.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetById retrieves a user by ID (alias for Get)
func (r *UserRepository) GetById(ctx context.Context, id int) (*models.User, error) {
	return r.Get(ctx, id)
}

// UpdatePassword updates the user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, hashedPassword string) error {
	result := r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword)
	return result.Error
}

// GetAll retrieves all users
func (r *UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	result := r.DB.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	result := r.DB.WithContext(ctx).Save(user)
	return result.Error
}

// Delete removes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	result := r.DB.WithContext(ctx).Delete(&models.User{}, id)
	return result.Error
}

// UpdateLastLogin updates the last login timestamp for a user
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id int) error {
	now := time.Now()
	result := r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("last_login", now)
	return result.Error
}

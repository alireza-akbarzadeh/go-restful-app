package repository

import (
	"context"
	"errors"
	"time"

	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/pagination"
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
	logging.Debug(ctx, "creating new user", "email", user.Email)

	result := r.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		logging.Error(ctx, "failed to create user", result.Error, "email", user.Email)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to create user")
	}

	logging.Info(ctx, "user created successfully", "user_id", user.ID, "email", user.Email)
	return user, nil
}

// Get retrieves a user by ID
func (r *UserRepository) Get(ctx context.Context, id int) (*models.User, error) {
	logging.Debug(ctx, "retrieving user by ID", "user_id", id)

	var user models.User
	result := r.DB.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logging.Debug(ctx, "user not found", "user_id", id)
			return nil, appErrors.Newf(appErrors.ErrNotFound, "user with ID %d not found", id)
		}
		logging.Error(ctx, "failed to retrieve user", result.Error, "user_id", id)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve user")
	}

	logging.Debug(ctx, "user retrieved successfully", "user_id", id)
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	logging.Debug(ctx, "retrieving user by email", "email", email)

	var user models.User
	result := r.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logging.Debug(ctx, "user not found by email", "email", email)
			return nil, appErrors.Newf(appErrors.ErrNotFound, "user with email %s not found", email)
		}
		logging.Error(ctx, "failed to retrieve user by email", result.Error, "email", email)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve user by email")
	}

	logging.Debug(ctx, "user retrieved successfully by email", "user_id", user.ID, "email", email)
	return &user, nil
}

// UpdatePassword updates the user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, hashedPassword string) error {
	logging.Debug(ctx, "updating user password", "user_id", userID)

	result := r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword)
	if result.Error != nil {
		logging.Error(ctx, "failed to update user password", result.Error, "user_id", userID)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to update password")
	}

	if result.RowsAffected == 0 {
		logging.Debug(ctx, "no user found to update password", "user_id", userID)
		return appErrors.Newf(appErrors.ErrNotFound, "user with ID %d not found", userID)
	}

	logging.Info(ctx, "user password updated successfully", "user_id", userID)
	return nil
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
	logging.Debug(ctx, "updating user", "user_id", user.ID, "email", user.Email)

	result := r.DB.WithContext(ctx).Save(user)
	if result.Error != nil {
		logging.Error(ctx, "failed to update user", result.Error, "user_id", user.ID)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to update user")
	}

	logging.Info(ctx, "user updated successfully", "user_id", user.ID, "email", user.Email)
	return nil
}

// Delete removes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	logging.Debug(ctx, "deleting user", "user_id", id)

	result := r.DB.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		logging.Error(ctx, "failed to delete user", result.Error, "user_id", id)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to delete user")
	}

	if result.RowsAffected == 0 {
		logging.Debug(ctx, "no user found to delete", "user_id", id)
		return appErrors.Newf(appErrors.ErrNotFound, "user with ID %d not found", id)
	}

	logging.Info(ctx, "user deleted successfully", "user_id", id)
	return nil
}

// UpdateLastLogin updates the last login timestamp for a user
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id int) error {
	logging.Debug(ctx, "updating last login", "user_id", id)

	now := time.Now()
	result := r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("last_login", now)
	if result.Error != nil {
		logging.Error(ctx, "failed to update last login", result.Error, "user_id", id)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to update last login")
	}

	if result.RowsAffected == 0 {
		logging.Debug(ctx, "no user found to update last login", "user_id", id)
		return appErrors.Newf(appErrors.ErrNotFound, "user with ID %d not found", id)
	}

	logging.Debug(ctx, "last login updated successfully", "user_id", id)
	return nil
}

// ListWithPagination retrieves users with pagination
func (r *UserRepository) ListWithPagination(ctx context.Context, req *pagination.PaginationRequest) ([]*models.User, *pagination.PaginationResponse, error) {
	logging.Debug(ctx, "retrieving users with pagination", "page", req.Page, "page_size", req.PageSize)

	var users []*models.User
	var total int64

	// Count total records
	if err := r.DB.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		logging.Error(ctx, "failed to count users", err)
		return nil, nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to count users")
	}

	// Get paginated records
	if err := r.DB.WithContext(ctx).
		Offset(req.Offset()).
		Limit(req.PageSize).
		Find(&users).Error; err != nil {
		logging.Error(ctx, "failed to retrieve users", err)
		return nil, nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve users")
	}

	// Calculate pagination response
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	paginationResp := &pagination.PaginationResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}

	logging.Info(ctx, "users retrieved successfully", "count", len(users), "total", total, "page", req.Page)
	return users, paginationResp, nil
}

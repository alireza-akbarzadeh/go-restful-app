package repository

import (
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Name     string `json:"name" gorm:"not null"`
	Password string `json:"-" gorm:"not null"` // Never expose password in JSON
}

// UserRepository handles user database operations
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Insert creates a new user in the database
func (r *UserRepository) Insert(user *User) (*User, error) {
	result := r.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// Get retrieves a user by ID
func (r *UserRepository) Get(id int) (*User, error) {
	var user User
	result := r.DB.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetById retrieves a user by ID (alias for Get)
func (r *UserRepository) GetById(id int) (*User, error) {
	return r.Get(id)
}

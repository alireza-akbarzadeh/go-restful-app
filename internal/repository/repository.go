package repository

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/repository/interfaces"
	"gorm.io/gorm"
)

// Models holds all repository models
type Models struct {
	Users      interfaces.UserRepositoryInterface
	Events     interfaces.EventRepositoryInterface
	Attendees  interfaces.AttendeeRepositoryInterface
	Categories interfaces.CategoryRepositoryInterface
	Comments   interfaces.CommentRepositoryInterface
	Profiles   interfaces.ProfileRepositoryInterface
	Products   interfaces.ProductRepositoryInterface
	Baskets    interfaces.BasketRepositoryInterface
	TxManager  *TxManager
}

// NewModels creates a new Models instance with all repositories
func NewModels(db *gorm.DB) *Models {
	return &Models{
		Users:      NewUserRepository(db),
		Events:     NewEventRepository(db),
		Attendees:  NewAttendeeRepository(db),
		Categories: NewCategoryRepository(db),
		Comments:   NewCommentRepository(db),
		Profiles:   NewProfileRepository(db),
		Products:   NewProductRepository(db),
		Baskets:    NewBasketRepository(db),
		TxManager:  NewTxManager(db),
	}
}

package repository

import "gorm.io/gorm"

// Models holds all repository models
type Models struct {
	Users      *UserRepository
	Events     *EventRepository
	Attendees  *AttendeeRepository
	Categories *CategoryRepository
	Comments   *CommentRepository
	Profiles   *ProfileRepository
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
	}
}

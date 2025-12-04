package repository

import "database/sql"

// Models holds all repository models
type Models struct {
	Users     *UserRepository
	Events    *EventRepository
	Attendees *AttendeeRepository
}

// NewModels creates a new Models instance with all repositories
func NewModels(db *sql.DB) *Models {
	return &Models{
		Users:     NewUserRepository(db),
		Events:    NewEventRepository(db),
		Attendees: NewAttendeeRepository(db),
	}
}

package repository

import (
	"gorm.io/gorm"
)

// Attendee represents an attendee relationship between a user and an event
type Attendee struct {
	ID      int `json:"id" gorm:"primaryKey"`
	UserID  int `json:"userId" gorm:"not null"`
	EventID int `json:"eventId" gorm:"not null"`
}

// AttendeeRepository handles attendee database operations
type AttendeeRepository struct {
	DB *gorm.DB
}

// NewAttendeeRepository creates a new AttendeeRepository
func NewAttendeeRepository(db *gorm.DB) *AttendeeRepository {
	return &AttendeeRepository{DB: db}
}

// Insert creates a new attendee record
func (r *AttendeeRepository) Insert(attendee *Attendee) (*Attendee, error) {
	result := r.DB.Create(attendee)
	if result.Error != nil {
		return nil, result.Error
	}
	return attendee, nil
}

// GetByEventAndUser retrieves an attendee record by event ID and user ID
func (r *AttendeeRepository) GetByEventAndUser(eventID, userID int) (*Attendee, error) {
	var attendee Attendee
	result := r.DB.Where("event_id = ? AND user_id = ?", eventID, userID).First(&attendee)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &attendee, nil
}

// GetByEventAndAttendee is an alias for GetByEventAndUser for backwards compatibility
func (r *AttendeeRepository) GetByEventAndAttendee(eventID, userID int) (*Attendee, error) {
	return r.GetByEventAndUser(eventID, userID)
}

// GetAttendeesByEvent retrieves all users attending a specific event
func (r *AttendeeRepository) GetAttendeesByEvent(eventID int) ([]*User, error) {
	var users []*User
	// Using a JOIN query to fetch users directly
	err := r.DB.Table("users").
		Joins("JOIN attendees ON attendees.user_id = users.id").
		Where("attendees.event_id = ?", eventID).
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetEventsByAttendee retrieves all events that a user is attending
func (r *AttendeeRepository) GetEventsByAttendee(userID int) ([]*Event, error) {
	var events []*Event
	// Using a JOIN query to fetch events directly
	err := r.DB.Table("events").
		Joins("JOIN attendees ON attendees.event_id = events.id").
		Where("attendees.user_id = ?", userID).
		Find(&events).Error

	if err != nil {
		return nil, err
	}
	return events, nil
}

// GetEventByAttendee is an alias for GetEventsByAttendee for backwards compatibility
func (r *AttendeeRepository) GetEventByAttendee(userID int) ([]*Event, error) {
	return r.GetEventsByAttendee(userID)
}

// Delete removes an attendee record
func (r *AttendeeRepository) Delete(userID, eventID int) error {
	result := r.DB.Where("user_id = ? AND event_id = ?", userID, eventID).Delete(&Attendee{})
	return result.Error
}

// DeleteByEvent removes all attendees for a specific event
func (r *AttendeeRepository) DeleteByEvent(eventID int) error {
	result := r.DB.Where("event_id = ?", eventID).Delete(&Attendee{})
	return result.Error
}

// DeleteByUser removes all attendee records for a specific user
func (r *AttendeeRepository) DeleteByUser(userID int) error {
	result := r.DB.Where("user_id = ?", userID).Delete(&Attendee{})
	return result.Error
}

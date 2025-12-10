package repository

import (
	"context"

	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"gorm.io/gorm"
)

// AttendeeRepository handles attendee database operations
type AttendeeRepository struct {
	DB *gorm.DB
}

// NewAttendeeRepository creates a new AttendeeRepository
func NewAttendeeRepository(db *gorm.DB) *AttendeeRepository {
	return &AttendeeRepository{DB: db}
}

// Insert creates a new attendee record
func (r *AttendeeRepository) Insert(ctx context.Context, attendee *models.Attendee) (*models.Attendee, error) {
	logging.Debug(ctx, "creating attendee registration", "event_id", attendee.EventID, "user_id", attendee.UserID)

	result := r.DB.WithContext(ctx).Create(attendee)
	if result.Error != nil {
		logging.Error(ctx, "failed to create attendee", result.Error, "event_id", attendee.EventID, "user_id", attendee.UserID)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to register attendee")
	}

	logging.Info(ctx, "attendee registered successfully", "attendee_id", attendee.ID, "event_id", attendee.EventID, "user_id", attendee.UserID)
	return attendee, nil
}

// GetByEventAndUser retrieves an attendee record by event ID and user ID
func (r *AttendeeRepository) GetByEventAndUser(ctx context.Context, eventID, userID int) (*models.Attendee, error) {
	var attendee models.Attendee
	result := r.DB.WithContext(ctx).Where("event_id = ? AND user_id = ?", eventID, userID).First(&attendee)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &attendee, nil
}

// GetByEventAndAttendee is an alias for GetByEventAndUser for backwards compatibility
func (r *AttendeeRepository) GetByEventAndAttendee(ctx context.Context, eventID, userID int) (*models.Attendee, error) {
	return r.GetByEventAndUser(ctx, eventID, userID)
}

// GetAttendeesByEvent retrieves all users attending a specific event
func (r *AttendeeRepository) GetAttendeesByEvent(ctx context.Context, eventID int) ([]*models.User, error) {
	var users []*models.User
	// Using a JOIN query to fetch users directly
	err := r.DB.WithContext(ctx).Table("users").
		Joins("JOIN attendees ON attendees.user_id = users.id").
		Where("attendees.event_id = ?", eventID).
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetEventsByAttendee retrieves all events that a user is attending
func (r *AttendeeRepository) GetEventsByAttendee(ctx context.Context, userID int) ([]*models.Event, error) {
	var events []*models.Event
	// Using a JOIN query to fetch events directly
	err := r.DB.WithContext(ctx).Table("events").
		Joins("JOIN attendees ON attendees.event_id = events.id").
		Where("attendees.user_id = ?", userID).
		Find(&events).Error

	if err != nil {
		return nil, err
	}
	return events, nil
}

// GetEventByAttendee is an alias for GetEventsByAttendee for backwards compatibility
func (r *AttendeeRepository) GetEventByAttendee(ctx context.Context, userID int) ([]*models.Event, error) {
	return r.GetEventsByAttendee(ctx, userID)
}

// Delete removes an attendee record
func (r *AttendeeRepository) Delete(ctx context.Context, userID, eventID int) error {
	result := r.DB.WithContext(ctx).Where("user_id = ? AND event_id = ?", userID, eventID).Delete(&models.Attendee{})
	return result.Error
}

// DeleteByEvent removes all attendees for a specific event
func (r *AttendeeRepository) DeleteByEvent(ctx context.Context, eventID int) error {
	result := r.DB.WithContext(ctx).Where("event_id = ?", eventID).Delete(&models.Attendee{})
	return result.Error
}

// DeleteByUser removes all attendee records for a specific user
func (r *AttendeeRepository) DeleteByUser(ctx context.Context, userID int) error {
	result := r.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.Attendee{})
	return result.Error
}

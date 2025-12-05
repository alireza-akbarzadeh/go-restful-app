package interfaces

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
)

type AttendeeRepositoryInterface interface {
	Insert(ctx context.Context, attendee *models.Attendee) (*models.Attendee, error)
	GetByEventAndUser(ctx context.Context, eventID, userID int) (*models.Attendee, error)
	GetByEventAndAttendee(ctx context.Context, eventID, userID int) (*models.Attendee, error)
	GetAttendeesByEvent(ctx context.Context, eventID int) ([]*models.User, error)
	GetEventsByAttendee(ctx context.Context, userID int) ([]*models.Event, error)
	GetEventByAttendee(ctx context.Context, userID int) ([]*models.Event, error)
	Delete(ctx context.Context, userID, eventID int) error
	DeleteByEvent(ctx context.Context, eventID int) error
	DeleteByUser(ctx context.Context, userID int) error
}

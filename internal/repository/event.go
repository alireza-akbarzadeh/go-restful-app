package repository

import (
	"context"
	"errors"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"gorm.io/gorm"
)

// EventRepository handles event database operations
type EventRepository struct {
	DB *gorm.DB
}

// NewEventRepository creates a new EventRepository
func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{DB: db}
}

// Insert creates a new event in the database
func (r *EventRepository) Insert(ctx context.Context, event *models.Event) (*models.Event, error) {
	result := r.DB.WithContext(ctx).Create(event)
	if result.Error != nil {
		return nil, result.Error
	}
	return event, nil
}

// Get retrieves an event by ID
func (r *EventRepository) Get(ctx context.Context, id int) (*models.Event, error) {
	var event models.Event
	result := r.DB.WithContext(ctx).First(&event, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &event, nil
}

// GetAll retrieves all events
func (r *EventRepository) GetAll(ctx context.Context) ([]*models.Event, error) {
	var events []*models.Event
	result := r.DB.WithContext(ctx).Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

// Update updates an existing event
func (r *EventRepository) Update(ctx context.Context, event *models.Event) error {
	result := r.DB.WithContext(ctx).Save(event)
	return result.Error
}

// Delete removes an event by ID
func (r *EventRepository) Delete(ctx context.Context, id int) error {
	result := r.DB.WithContext(ctx).Delete(&models.Event{}, id)
	return result.Error
}

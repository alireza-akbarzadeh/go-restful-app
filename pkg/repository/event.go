package repository

import (
	"errors"

	"github.com/alireza-akbarzadeh/ginflow/pkg/models"
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
func (r *EventRepository) Insert(event *models.Event) (*models.Event, error) {
	result := r.DB.Create(event)
	if result.Error != nil {
		return nil, result.Error
	}
	return event, nil
}

// Get retrieves an event by ID
func (r *EventRepository) Get(id int) (*models.Event, error) {
	var event models.Event
	result := r.DB.First(&event, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &event, nil
}

// GetAll retrieves all events
func (r *EventRepository) GetAll() ([]*models.Event, error) {
	var events []*models.Event
	result := r.DB.Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

// Update updates an existing event
func (r *EventRepository) Update(event *models.Event) error {
	result := r.DB.Save(event)
	return result.Error
}

// Delete removes an event by ID
func (r *EventRepository) Delete(id int) error {
	result := r.DB.Delete(&models.Event{}, id)
	return result.Error
}

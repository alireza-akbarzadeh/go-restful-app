package repository

import (
	"gorm.io/gorm"
)

// Event represents an event in the system
type Event struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	OwnerID     int    `json:"ownerId" gorm:"not null"`
	Name        string `json:"name" binding:"required,min=3" gorm:"not null"`
	Description string `json:"description" binding:"required,min=10" gorm:"not null"`
	Date        string `json:"date" binding:"required,datetime=2006-01-02" gorm:"not null"`
	Location    string `json:"location" binding:"required,min=3" gorm:"not null"`
}

// EventRepository handles event database operations
type EventRepository struct {
	DB *gorm.DB
}

// NewEventRepository creates a new EventRepository
func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{DB: db}
}

// Insert creates a new event in the database
func (r *EventRepository) Insert(event *Event) (*Event, error) {
	result := r.DB.Create(event)
	if result.Error != nil {
		return nil, result.Error
	}
	return event, nil
}

// Get retrieves an event by ID
func (r *EventRepository) Get(id int) (*Event, error) {
	var event Event
	result := r.DB.First(&event, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &event, nil
}

// GetAll retrieves all events
func (r *EventRepository) GetAll() ([]*Event, error) {
	var events []*Event
	result := r.DB.Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

// Update updates an existing event
func (r *EventRepository) Update(event *Event) error {
	result := r.DB.Save(event)
	return result.Error
}

// Delete removes an event by ID
func (r *EventRepository) Delete(id int) error {
	result := r.DB.Delete(&Event{}, id)
	return result.Error
}

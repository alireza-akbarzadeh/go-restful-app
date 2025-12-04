package repository

import (
	"context"
	"database/sql"
	"time"
)

// Event represents an event in the system
type Event struct {
	ID          int    `json:"id"`
	OwnerID     int    `json:"ownerId"`
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	Date        string `json:"date" binding:"required,datetime=2006-01-02"`
	Location    string `json:"location" binding:"required,min=3"`
}

// EventRepository handles event database operations
type EventRepository struct {
	DB *sql.DB
}

// NewEventRepository creates a new EventRepository
func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{DB: db}
}

// Insert creates a new event in the database
func (r *EventRepository) Insert(event *Event) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO events (owner_id, name, description, date, location) VALUES (?, ?, ?, ?, ?)`
	result, err := r.DB.ExecContext(ctx, query,
		event.OwnerID, event.Name, event.Description, event.Date, event.Location)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	event.ID = int(id)

	return event, nil
}

// Get retrieves an event by ID
func (r *EventRepository) Get(id int) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, owner_id, name, description, date, location FROM events WHERE id = ?`
	var event Event
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.OwnerID,
		&event.Name,
		&event.Description,
		&event.Date,
		&event.Location,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}

// GetAll retrieves all events
func (r *EventRepository) GetAll() ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, owner_id, name, description, date, location FROM events`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*Event{}
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.ID, &event.OwnerID, &event.Name, &event.Description, &event.Date, &event.Location); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// GetByOwner retrieves all events created by a specific owner
func (r *EventRepository) GetByOwner(ownerID int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, owner_id, name, description, date, location FROM events WHERE owner_id = ?`
	rows, err := r.DB.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*Event{}
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.ID, &event.OwnerID, &event.Name, &event.Description, &event.Date, &event.Location); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// Update updates an event's information
func (r *EventRepository) Update(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE events SET owner_id = ?, name = ?, description = ?, date = ?, location = ? WHERE id = ?`
	_, err := r.DB.ExecContext(ctx, query,
		event.OwnerID, event.Name, event.Description, event.Date, event.Location, event.ID)
	return err
}

// Delete deletes an event by ID
func (r *EventRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM events WHERE id = ?`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}

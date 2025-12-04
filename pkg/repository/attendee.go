package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Attendee represents an attendee relationship between a user and an event
type Attendee struct {
	ID      int `json:"id"`
	UserID  int `json:"userId"`
	EventID int `json:"eventId"`
}

// AttendeeRepository handles attendee database operations
type AttendeeRepository struct {
	DB *sql.DB
}

// NewAttendeeRepository creates a new AttendeeRepository
func NewAttendeeRepository(db *sql.DB) *AttendeeRepository {
	return &AttendeeRepository{DB: db}
}

// Insert creates a new attendee record
func (r *AttendeeRepository) Insert(attendee *Attendee) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO attendees (event_id, user_id) VALUES (?, ?)`
	result, err := r.DB.ExecContext(ctx, query, attendee.EventID, attendee.UserID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	attendee.ID = int(id)

	return attendee, nil
}

// GetByEventAndUser retrieves an attendee record by event ID and user ID
func (r *AttendeeRepository) GetByEventAndUser(eventID, userID int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, event_id FROM attendees WHERE event_id = ? AND user_id = ?`
	var attendee Attendee
	err := r.DB.QueryRowContext(ctx, query, eventID, userID).Scan(
		&attendee.ID,
		&attendee.UserID,
		&attendee.EventID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &attendee, nil
}

// GetByEventAndAttendee is an alias for GetByEventAndUser for backwards compatibility
func (r *AttendeeRepository) GetByEventAndAttendee(eventID, userID int) (*Attendee, error) {
	return r.GetByEventAndUser(eventID, userID)
}

// GetAttendeesByEvent retrieves all users attending a specific event
func (r *AttendeeRepository) GetAttendeesByEvent(eventID int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT u.id, u.name, u.email
		FROM users u
		JOIN attendees a ON u.id = a.user_id
		WHERE a.event_id = ?
	`
	rows, err := r.DB.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetEventsByAttendee retrieves all events that a user is attending
func (r *AttendeeRepository) GetEventsByAttendee(userID int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT e.id, e.owner_id, e.name, e.description, e.date, e.location
		FROM events e
		JOIN attendees a ON e.id = a.event_id
		WHERE a.user_id = ?
	`
	rows, err := r.DB.QueryContext(ctx, query, userID)
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

// GetEventByAttendee is an alias for GetEventsByAttendee for backwards compatibility
func (r *AttendeeRepository) GetEventByAttendee(userID int) ([]*Event, error) {
	return r.GetEventsByAttendee(userID)
}

// Delete removes an attendee record
func (r *AttendeeRepository) Delete(userID, eventID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM attendees WHERE user_id = ? AND event_id = ?`
	_, err := r.DB.ExecContext(ctx, query, userID, eventID)
	return err
}

// DeleteByEvent removes all attendees for a specific event
func (r *AttendeeRepository) DeleteByEvent(eventID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM attendees WHERE event_id = ?`
	_, err := r.DB.ExecContext(ctx, query, eventID)
	return err
}

// DeleteByUser removes all attendee records for a specific user
func (r *AttendeeRepository) DeleteByUser(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM attendees WHERE user_id = ?`
	_, err := r.DB.ExecContext(ctx, query, userID)
	return err
}

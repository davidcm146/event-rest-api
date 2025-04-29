package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	Id      string `json:"id"`
	UserId  string `json:"userId"`
	EventId string `json:"eventId"`
}

func (m *AttendeeModel) Insert(attendee *Attendee) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `INSERT INTO attendees (user_id, event_id) VALUES ($1, $2) RETURNING id`
	err := m.DB.QueryRowContext(ctx, query, attendee.UserId, attendee.EventId).Scan(&attendee.Id)
	if err != nil {
		return nil, err
	}
	return attendee, nil
}

func (m *AttendeeModel) GetByEventAndAttendee(eventId, userId string) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, user_id, event_id FROM attendees WHERE event_id = $1 AND user_id = $2`
	var attendee Attendee
	err := m.DB.QueryRowContext(ctx, query, eventId, userId).Scan(&attendee.Id, &attendee.UserId, &attendee.EventId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No attendee found
		}
		return nil, err // Some other error occurred
	}
	return &attendee, nil // Attendee found
}

func (m *AttendeeModel) GetAttendeesByEventId(eventId string) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT u.id, u.name, u.email FROM attendees a JOIN users u ON a.user_id = u.id WHERE event_id = $1`
	var attendees []*User
	rows, err := m.DB.QueryContext(ctx, query, eventId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		row := &User{}
		if err := rows.Scan(&row.Id, &row.Name, &row.Email); err != nil {
			return nil, err
		}
		attendees = append(attendees, row)
	}
	return attendees, nil
}

func (m *AttendeeModel) GetEventsByAttendeeId(userId string) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT e.id, e.name, e.owner_id, e.description, e.date, e.location FROM attendees a JOIN events e ON a.event_id = e.id WHERE a.user_id = $1`
	var events []*Event
	rows, err := m.DB.QueryContext(ctx, query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		row := &Event{}
		if err := rows.Scan(&row.Id, &row.Name, &row.OwnerId, &row.Description, &row.Date, &row.Location); err != nil {
			return nil, err
		}
		events = append(events, row)
	}
	return events, nil
}

func (m *AttendeeModel) Delete(userId, eventId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM attendees WHERE user_id = $1 AND event_id = $2`
	_, err := m.DB.ExecContext(ctx, query, userId, eventId)
	if err != nil {
		return err
	}
	return nil
}

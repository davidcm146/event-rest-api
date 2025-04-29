package database

import (
	"context"
	"database/sql"
	"github.com/davidcm146/event-rest-api/internal/utils"
	"time"
)

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	Id          string `json:"id"`
	Name        string `json:"name" binding:"required,min=3"`
	OwnerId     string `json:"ownerId" binding:"required"`
	Description string `json:"description" binding:"required,min=10"`
	Date        string `json:"date" binding:"required,datetime=02/01/2006"`
	Location    string `json:"location"`
}

func (m *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	formattedDate, err := utils.ParseAndFormatDate(event.Date)
	if err != nil {
		return err
	}
	query := `INSERT INTO events (name, owner_id, description, date, location) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return m.DB.QueryRowContext(ctx, query, event.Name, event.OwnerId, event.Description, formattedDate, event.Location).Scan(&event.Id)
}

func (m *EventModel) GetAll() ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT id, name, owner_id, description, date, location FROM events`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*Event{}
	for rows.Next() {
		event := &Event{}
		if err := rows.Scan(&event.Id, &event.Name, &event.OwnerId, &event.Description, &event.Date, &event.Location); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (m *EventModel) Get(id string) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT id, name, owner_id, description, date, location FROM events WHERE id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)

	event := Event{}
	if err := row.Scan(&event.Id, &event.Name, &event.OwnerId, &event.Description, &event.Date, &event.Location); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (m *EventModel) Update(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	formattedDate, err := utils.ParseAndFormatDate(event.Date)
	if err != nil {
		return err
	}
	query := `UPDATE events SET name = $1, description = $2, date = $3, location = $4 WHERE id = $5`
	_, err = m.DB.ExecContext(ctx, query, event.Name, event.Description, formattedDate, event.Location, event.Id)
	if err != nil {
		return err
	}
	return nil
}

func (m *EventModel) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM events WHERE id = $1`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

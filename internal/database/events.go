package database

import "database/sql"

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	Id          string `json:"id"`
	Name        string `json:"name" binding:"required, min=3"`
	OwnerId     string `json:"ownerId" binding:"required"`
	Description string `json:"description" binding:"required, min=10"`
	Date        string `json:"date" binding:"required, datetime=20/04/2025"`
	Location    string `json:"location"`
}

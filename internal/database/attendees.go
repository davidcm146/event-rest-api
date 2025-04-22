package database

import "database/sql"

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	Id      string `json:"id"`
	UserId  string `json:"userId"`
	EventId string `json:"eventId"`
}

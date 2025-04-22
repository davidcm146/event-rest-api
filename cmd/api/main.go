package main

import (
	"database/sql"
	"log"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:10042003@localhost:5432/events_api?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}

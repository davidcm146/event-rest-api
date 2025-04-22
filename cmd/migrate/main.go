package main

import (
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a migration command up or down")
	}

	direction := os.Args[1]

	db, err := sql.Open("postgres", "postgres://postgres:10042003@localhost:5432/events_api?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "postgres", instance)
	if err != nil {
		log.Fatal(err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	default:
		log.Fatal("Invalid command. Use 'up' or 'down'.")
	}
}

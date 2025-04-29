package main

import (
	"database/sql"
	_ "github.com/davidcm146/event-rest-api/docs"
	"github.com/davidcm146/event-rest-api/internal/database"
	"github.com/davidcm146/event-rest-api/internal/env"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
)

// @title Event REST API
// @version 1.0
// @description This is a simple REST API for managing events
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your Bearer token in the format **Bearer &lt;token&gt;**
type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:10042003@localhost:5432/events_api?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "defaultsecret"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}

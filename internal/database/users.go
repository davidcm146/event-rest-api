package database

import "database/sql"

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

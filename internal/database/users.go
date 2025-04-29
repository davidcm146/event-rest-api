package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the user already exists
	existingUser := &User{}
	checkQuery := `SELECT id, name, email FROM users WHERE email = $1`
	err := m.DB.QueryRowContext(ctx, checkQuery, user.Email).Scan(&existingUser.Id, &existingUser.Name, &existingUser.Email)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if err == nil {
		// Found an existing user
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	insertQuery := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	err = m.DB.QueryRowContext(ctx, insertQuery, user.Name, user.Email, user.Password).Scan(&user.Id)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (m *UserModel) GetUser(query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, args...)

	user := &User{}
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (m *UserModel) GetById(id string) (*User, error) {
	query := `SELECT id, name, email, password FROM users WHERE id = $1`
	return m.GetUser(query, id)
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, name, email, password FROM users WHERE email = $1`
	return m.GetUser(query, email)
}

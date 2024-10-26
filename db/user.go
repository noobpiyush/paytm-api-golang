package db

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("User already exists")
)

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	FullName     string `json:"full_name"`
}

func CreateUser(email, passwordHash, fullName string) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	query := `
        INSERT INTO users (email, password_hash, full_name) 
        VALUES ($1, $2, $3)`

	_, err := DB.Exec(query, email, passwordHash, fullName)
	if err != nil {
		// Check for unique constraint violation
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return ErrUserExists
		}
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(email string) (*User, error) {
	if DB == nil {
		return nil, errors.New("database not initialized")
	}

	var user User
	query := `
        SELECT id, email, password_hash, full_name 
        FROM users 
        WHERE email = $1`

	err := DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	return &user, nil
}

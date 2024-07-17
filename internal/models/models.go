// Package models consists of data models
package models

import "time"

// User represents the application user
type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	PassHash  []byte    `json:"-" db:"hashed_password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

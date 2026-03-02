package models

import "time"

// Role represents user role in the system
const (
	RolePassenger = "passenger"
	RoleDriver    = "driver"
)

// User represents a user in the system (passenger or driver)
type User struct {
	ID        string    `json:"id" db:"id"`
	Phone     string    `json:"phone" db:"phone"`
	Name      string    `json:"name" db:"name"`
	Role      string    `json:"role" db:"role"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

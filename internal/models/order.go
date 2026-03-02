package models

import "time"

// Order status constants
const (
	OrderStatusCreated   = "created"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)

// Order represents a taxi order (archive/history)
type Order struct {
	ID          string     `json:"id" db:"id"`
	PassengerID string     `json:"passenger_id" db:"passenger_id"`
	DriverID    *string    `json:"driver_id,omitempty" db:"driver_id"`
	Status      string     `json:"status" db:"status"`
	Price       int        `json:"price" db:"price"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

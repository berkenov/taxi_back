package models

// Car represents a driver's vehicle
type Car struct {
	ID       string `json:"id" db:"id"`
	DriverID string `json:"driver_id" db:"driver_id"`
	Model    string `json:"model" db:"model"`
	Number   string `json:"number" db:"number"`
	Color    string `json:"color" db:"color"`
}

package repository

import (
	"database/sql"
	"taxi/internal/models"

	"github.com/google/uuid"
)

// CarRepository handles car persistence
type CarRepository struct {
	db *sql.DB
}

// NewCarRepository creates a new CarRepository
func NewCarRepository(db *sql.DB) *CarRepository {
	return &CarRepository{db: db}
}

// Create inserts a new car
func (r *CarRepository) Create(car *models.Car) error {
	car.ID = uuid.New().String()
	query := `INSERT INTO cars (id, driver_id, model, number, color) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, car.ID, car.DriverID, car.Model, car.Number, car.Color)
	return err
}

// GetByID finds car by ID
func (r *CarRepository) GetByID(id string) (*models.Car, error) {
	query := `SELECT id, driver_id, model, number, color FROM cars WHERE id = $1`
	var c models.Car
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.DriverID, &c.Model, &c.Number, &c.Color)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetByDriverID returns all cars for a driver
func (r *CarRepository) GetByDriverID(driverID string) ([]*models.Car, error) {
	query := `SELECT id, driver_id, model, number, color FROM cars WHERE driver_id = $1 ORDER BY model`
	rows, err := r.db.Query(query, driverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []*models.Car
	for rows.Next() {
		var c models.Car
		if err := rows.Scan(&c.ID, &c.DriverID, &c.Model, &c.Number, &c.Color); err != nil {
			return nil, err
		}
		cars = append(cars, &c)
	}
	return cars, rows.Err()
}

// Update updates car by ID
func (r *CarRepository) Update(car *models.Car) error {
	query := `UPDATE cars SET model = $1, number = $2, color = $3 WHERE id = $4`
	result, err := r.db.Exec(query, car.Model, car.Number, car.Color, car.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Delete removes car by ID
func (r *CarRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM cars WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

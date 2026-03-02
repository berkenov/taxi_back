package service

import (
	"database/sql"
	"errors"

	"taxi/internal/models"
	"taxi/internal/repository"
)

var (
	ErrCarNotFound   = errors.New("car not found")
	ErrNotDriver     = errors.New("only drivers can have cars")
	ErrModelRequired = errors.New("model is required")
	ErrNumberRequired = errors.New("number is required")
	ErrColorRequired = errors.New("color is required")
)

// CarService contains business logic for cars
type CarService struct {
	carRepo  *repository.CarRepository
	userRepo *repository.UserRepository
}

// NewCarService creates a new CarService
func NewCarService(carRepo *repository.CarRepository, userRepo *repository.UserRepository) *CarService {
	return &CarService{carRepo: carRepo, userRepo: userRepo}
}

// Create adds a car for a driver (only drivers can have cars)
func (s *CarService) Create(driverID, model, number, color string) (*models.Car, error) {
	driver, err := s.userRepo.GetByID(driverID)
	if err != nil || driver == nil {
		return nil, ErrUserNotFound
	}
	if driver.Role != models.RoleDriver {
		return nil, ErrNotDriver
	}
	if model == "" {
		return nil, ErrModelRequired
	}
	if number == "" {
		return nil, ErrNumberRequired
	}
	if color == "" {
		return nil, ErrColorRequired
	}

	car := &models.Car{
		DriverID: driverID,
		Model:    model,
		Number:   number,
		Color:    color,
	}
	if err := s.carRepo.Create(car); err != nil {
		return nil, err
	}
	return car, nil
}

// GetByID returns car by ID
func (s *CarService) GetByID(id string) (*models.Car, error) {
	car, err := s.carRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if car == nil {
		return nil, ErrCarNotFound
	}
	return car, nil
}

// GetByDriverID returns all cars for a driver
func (s *CarService) GetByDriverID(driverID string) ([]*models.Car, error) {
	driver, err := s.userRepo.GetByID(driverID)
	if err != nil || driver == nil {
		return nil, ErrUserNotFound
	}
	if driver.Role != models.RoleDriver {
		return nil, ErrNotDriver
	}
	return s.carRepo.GetByDriverID(driverID)
}

// Update updates car by ID
func (s *CarService) Update(id string, model, number, color *string) (*models.Car, error) {
	car, err := s.carRepo.GetByID(id)
	if err != nil || car == nil {
		return nil, ErrCarNotFound
	}

	if model != nil {
		if *model == "" {
			return nil, ErrModelRequired
		}
		car.Model = *model
	}
	if number != nil {
		if *number == "" {
			return nil, ErrNumberRequired
		}
		car.Number = *number
	}
	if color != nil {
		if *color == "" {
			return nil, ErrColorRequired
		}
		car.Color = *color
	}

	if err := s.carRepo.Update(car); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCarNotFound
		}
		return nil, err
	}
	return car, nil
}

// Delete removes car by ID
func (s *CarService) Delete(id string) error {
	err := s.carRepo.Delete(id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return ErrCarNotFound
	}
	return err
}

package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"taxi/internal/models"
	"taxi/internal/service"
)

// CarHandler handles car CRUD endpoints (only for drivers)
type CarHandler struct {
	carSvc *service.CarService
}

// NewCarHandler creates a new CarHandler
func NewCarHandler(carSvc *service.CarService) *CarHandler {
	return &CarHandler{carSvc: carSvc}
}

// CreateRequest represents create car request body
type CreateCarRequest struct {
	DriverID string `json:"driver_id"`
	Model    string `json:"model"`
	Number   string `json:"number"`
	Color    string `json:"color"`
}

// UpdateCarRequest represents update car request body (all fields optional)
type UpdateCarRequest struct {
	Model  *string `json:"model,omitempty"`
	Number *string `json:"number,omitempty"`
	Color  *string `json:"color,omitempty"`
}

// Create adds a car for a driver
func (h *CarHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateCarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	car, err := h.carSvc.Create(req.DriverID, req.Model, req.Number, req.Color)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			http.Error(w, "Driver not found", http.StatusNotFound)
		case errors.Is(err, service.ErrNotDriver):
			http.Error(w, "Only drivers can have cars", http.StatusBadRequest)
		case errors.Is(err, service.ErrModelRequired):
			http.Error(w, "Model is required", http.StatusBadRequest)
		case errors.Is(err, service.ErrNumberRequired):
			http.Error(w, "Number is required", http.StatusBadRequest)
		case errors.Is(err, service.ErrColorRequired):
			http.Error(w, "Color is required", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(car)
}

// GetByID returns car by ID
func (h *CarHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Car ID required", http.StatusBadRequest)
		return
	}

	car, err := h.carSvc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrCarNotFound) {
			http.Error(w, "Car not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}

// ListByDriver returns all cars for a driver
func (h *CarHandler) ListByDriver(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	driverID := chi.URLParam(r, "driver_id")
	if driverID == "" {
		http.Error(w, "Driver ID required", http.StatusBadRequest)
		return
	}

	cars, err := h.carSvc.GetByDriverID(driverID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			http.Error(w, "Driver not found", http.StatusNotFound)
		case errors.Is(err, service.ErrNotDriver):
			http.Error(w, "Only drivers can have cars", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if cars == nil {
		cars = []*models.Car{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

// Update updates car by ID
func (h *CarHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Car ID required", http.StatusBadRequest)
		return
	}

	var req UpdateCarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	car, err := h.carSvc.Update(id, req.Model, req.Number, req.Color)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCarNotFound):
			http.Error(w, "Car not found", http.StatusNotFound)
		case errors.Is(err, service.ErrModelRequired):
			http.Error(w, "Model is required", http.StatusBadRequest)
		case errors.Is(err, service.ErrNumberRequired):
			http.Error(w, "Number is required", http.StatusBadRequest)
		case errors.Is(err, service.ErrColorRequired):
			http.Error(w, "Color is required", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}

// Delete removes car by ID
func (h *CarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Car ID required", http.StatusBadRequest)
		return
	}

	err := h.carSvc.Delete(id)
	if err != nil {
		if errors.Is(err, service.ErrCarNotFound) {
			http.Error(w, "Car not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

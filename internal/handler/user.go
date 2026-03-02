package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"taxi/internal/service"
)

// UserHandler handles user CRUD endpoints
type UserHandler struct {
	userSvc *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// RegisterRequest represents registration request body
type RegisterRequest struct {
	Phone string `json:"phone"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// UpdateRequest represents update request body (all fields optional)
type UpdateRequest struct {
	Phone   *string `json:"phone,omitempty"`
	Name    *string `json:"name,omitempty"`
	Role    *string `json:"role,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// Register creates a new user (passenger or driver)
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userSvc.Register(req.Phone, req.Name, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPhone):
			http.Error(w, "Invalid phone format (expected 7XXXXXXXXXX)", http.StatusBadRequest)
		case errors.Is(err, service.ErrNameRequired):
			http.Error(w, "Name is required", http.StatusBadRequest)
		case errors.Is(err, service.ErrInvalidRole):
			http.Error(w, "Role must be passenger or driver", http.StatusBadRequest)
		case errors.Is(err, service.ErrPhoneExists):
			http.Error(w, "Phone already registered", http.StatusConflict)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetByPhone returns user by phone (для входа существующего пользователя)
func (h *UserHandler) GetByPhone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	phone := chi.URLParam(r, "phone")
	if phone == "" {
		http.Error(w, "Phone required", http.StatusBadRequest)
		return
	}

	user, err := h.userSvc.GetUserByPhone(phone)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetByID returns user by ID
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	user, err := h.userSvc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Update updates user by ID
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userSvc.Update(id, req.Phone, req.Name, req.Role, req.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			http.Error(w, "User not found", http.StatusNotFound)
		case errors.Is(err, service.ErrInvalidPhone):
			http.Error(w, "Invalid phone format (expected 7XXXXXXXXXX)", http.StatusBadRequest)
		case errors.Is(err, service.ErrNameRequired):
			http.Error(w, "Name is required", http.StatusBadRequest)
		case errors.Is(err, service.ErrInvalidRole):
			http.Error(w, "Role must be passenger or driver", http.StatusBadRequest)
		case errors.Is(err, service.ErrPhoneExists):
			http.Error(w, "Phone already registered", http.StatusConflict)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Delete removes user by ID
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	err := h.userSvc.Delete(id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

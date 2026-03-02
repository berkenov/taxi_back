package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"taxi/internal/service"
)

// AuthHandler handles auth endpoints (OTP login)
type AuthHandler struct {
	authSvc *service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// AuthSendOTPRequest represents send OTP request
type AuthSendOTPRequest struct {
	Phone string `json:"phone"`
}

// AuthVerifyOTPRequest represents verify OTP request
type AuthVerifyOTPRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

// SendOTP generates code, sends via WhatsApp, stores for verification
func (h *AuthHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthSendOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.authSvc.SendOTP(req.Phone)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPhone) {
			http.Error(w, "Invalid phone format (expected 7XXXXXXXXXX)", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}

// VerifyOTP verifies code and returns user if exists
func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthVerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Code) != 4 {
		http.Error(w, "Code must be 4 digits", http.StatusBadRequest)
		return
	}

	user, err := h.authSvc.VerifyOTP(req.Phone, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPhone):
			http.Error(w, "Invalid phone format (expected 7XXXXXXXXXX)", http.StatusBadRequest)
		case errors.Is(err, service.ErrInvalidOTP):
			http.Error(w, "Invalid or expired code", http.StatusUnauthorized)
		case errors.Is(err, service.ErrUserNotFound):
			http.Error(w, "User not found. Please register first", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

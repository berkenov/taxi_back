package handler

import (
	"encoding/json"
	"net/http"

	"taxi/pkg/validator"
	"taxi/pkg/whatsapp"
)

// OTPHandler handles OTP-related endpoints
type OTPHandler struct {
	messenger whatsapp.Messenger
}

// NewOTPHandler creates a new OTPHandler
func NewOTPHandler(messenger whatsapp.Messenger) *OTPHandler {
	return &OTPHandler{messenger: messenger}
}

// SendOTPRequest represents the request body for sending OTP
type SendOTPRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

// SendOTP sends OTP to the specified phone (для проверки логирования Green API)
func (h *OTPHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !validator.ValidatePhone(req.Phone) {
		http.Error(w, "Invalid phone format (expected 7XXXXXXXXXX)", http.StatusBadRequest)
		return
	}

	phone := validator.NormalizePhone(req.Phone)
	if len(req.Code) != 4 {
		http.Error(w, "Code must be 4 digits", http.StatusBadRequest)
		return
	}

	if err := h.messenger.SendOTP(phone, req.Code); err != nil {
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}

package service

import (
	"crypto/rand"
	"errors"
	"fmt"

	"taxi/internal/models"
	"taxi/internal/repository"
	"taxi/pkg/validator"
	"taxi/pkg/whatsapp"
)

var ErrInvalidOTP = errors.New("invalid or expired OTP")

// AuthService handles authentication (OTP flow)
type AuthService struct {
	userRepo   *repository.UserRepository
	otpRepo    *repository.OTPRepository
	messenger  whatsapp.Messenger
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo *repository.UserRepository, otpRepo *repository.OTPRepository, messenger whatsapp.Messenger) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		otpRepo:   otpRepo,
		messenger: messenger,
	}
}

// generateCode returns random 4-digit code (0000-9999)
func generateCode() string {
	b := make([]byte, 2)
	rand.Read(b)
	n := int(b[0])<<8 | int(b[1])
	if n < 0 {
		n = -n
	}
	return fmt.Sprintf("%04d", n%10000)
}

// SendOTP generates code, sends via WhatsApp, stores in DB
func (s *AuthService) SendOTP(phone string) error {
	if !validator.ValidatePhone(phone) {
		return ErrInvalidPhone
	}
	normalized := validator.NormalizePhone(phone)
	code := generateCode()

	if err := s.messenger.SendOTP(normalized, code); err != nil {
		return err
	}
	return s.otpRepo.Store(normalized, code)
}

// VerifyOTP verifies code and returns user if exists
func (s *AuthService) VerifyOTP(phone, code string) (*models.User, error) {
	if !validator.ValidatePhone(phone) {
		return nil, ErrInvalidPhone
	}
	normalized := validator.NormalizePhone(phone)

	valid, err := s.otpRepo.Verify(normalized, code)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, ErrInvalidOTP
	}

	// Delete OTP after successful verification
	_ = s.otpRepo.Delete(normalized)

	user, err := s.userRepo.GetByPhone(normalized)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound // from user_service
	}

	return user, nil
}

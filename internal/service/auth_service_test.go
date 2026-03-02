package service

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"taxi/internal/repository"
)

type mockMessenger struct {
	sendErr error
}

func (m *mockMessenger) SendOTP(phone, code string) error {
	return m.sendErr
}

func TestAuthService_VerifyOTP_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	messenger := &mockMessenger{}
	svc := NewAuthService(userRepo, otpRepo, messenger)

	// OTP verify returns valid
	expiresAt := time.Now().Add(10 * time.Minute)
	mock.ExpectQuery("SELECT .+ FROM otp_codes").
		WithArgs("77001234567").
		WillReturnRows(sqlmock.NewRows([]string{"code", "expires_at"}).AddRow("1234", expiresAt))

	mock.ExpectExec("DELETE FROM otp_codes").
		WithArgs("77001234567").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// User exists
	mock.ExpectQuery("SELECT .+ FROM users WHERE phone").
		WithArgs("77001234567").
		WillReturnRows(sqlmock.NewRows([]string{"id", "phone", "name", "role", "is_active", "created_at"}).
			AddRow("user-id", "77001234567", "Test", "passenger", false, time.Now()))

	user, err := svc.VerifyOTP("77001234567", "1234")
	if err != nil {
		t.Fatalf("VerifyOTP: %v", err)
	}
	if user.Phone != "77001234567" {
		t.Errorf("phone = %q", user.Phone)
	}
}

func TestAuthService_VerifyOTP_InvalidCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	otpRepo := repository.NewOTPRepository(db)
	userRepo := repository.NewUserRepository(db)
	svc := NewAuthService(userRepo, otpRepo, &mockMessenger{})

	expiresAt := time.Now().Add(10 * time.Minute)
	mock.ExpectQuery("SELECT .+ FROM otp_codes").
		WithArgs("77001234567").
		WillReturnRows(sqlmock.NewRows([]string{"code", "expires_at"}).AddRow("9999", expiresAt))

	_, err = svc.VerifyOTP("77001234567", "1234")
	if err != ErrInvalidOTP {
		t.Errorf("VerifyOTP: want ErrInvalidOTP, got %v", err)
	}
}

func TestAuthService_VerifyOTP_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	otpRepo := repository.NewOTPRepository(db)
	userRepo := repository.NewUserRepository(db)
	svc := NewAuthService(userRepo, otpRepo, &mockMessenger{})

	expiresAt := time.Now().Add(10 * time.Minute)
	mock.ExpectQuery("SELECT .+ FROM otp_codes").
		WithArgs("77001234567").
		WillReturnRows(sqlmock.NewRows([]string{"code", "expires_at"}).AddRow("1234", expiresAt))

	mock.ExpectExec("DELETE FROM otp_codes").
		WithArgs("77001234567").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery("SELECT .+ FROM users WHERE phone").
		WithArgs("77001234567").
		WillReturnError(sql.ErrNoRows)

	_, err = svc.VerifyOTP("77001234567", "1234")
	if err != ErrUserNotFound {
		t.Errorf("VerifyOTP: want ErrUserNotFound, got %v", err)
	}
}

func TestAuthService_SendOTP_InvalidPhone(t *testing.T) {
	svc := NewAuthService(nil, nil, &mockMessenger{})
	err := svc.SendOTP("123")
	if err != ErrInvalidPhone {
		t.Errorf("SendOTP: want ErrInvalidPhone, got %v", err)
	}
}

func TestAuthService_SendOTP_MessengerError(t *testing.T) {
	svc := NewAuthService(nil, nil, &mockMessenger{sendErr: errors.New("whatsapp error")})
	err := svc.SendOTP("77001234567")
	if err == nil {
		t.Error("SendOTP: expected error")
	}
}

package handler

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"taxi/internal/repository"
	"taxi/internal/service"
)

type mockMessenger struct{ err error }

func (m *mockMessenger) SendOTP(phone, code string) error { return m.err }

func TestAuthHandler_SendOTP_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	messenger := &mockMessenger{}
	authSvc := service.NewAuthService(userRepo, otpRepo, messenger)

	mock.ExpectExec("INSERT INTO otp_codes").
		WithArgs("77001234567", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	h := NewAuthHandler(authSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/send-otp", bytes.NewReader([]byte(`{"phone":"77001234567"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.SendOTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestAuthHandler_SendOTP_InvalidPhone(t *testing.T) {
	authSvc := service.NewAuthService(nil, nil, &mockMessenger{})
	h := NewAuthHandler(authSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/send-otp", bytes.NewReader([]byte(`{"phone":"123"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.SendOTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestAuthHandler_VerifyOTP_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	authSvc := service.NewAuthService(userRepo, otpRepo, &mockMessenger{})

	expiresAt := time.Now().Add(10 * time.Minute)
	mock.ExpectQuery("SELECT .+ FROM otp_codes").
		WithArgs("77001234567").
		WillReturnRows(sqlmock.NewRows([]string{"code", "expires_at"}).AddRow("1234", expiresAt))

	mock.ExpectExec("DELETE FROM otp_codes").
		WithArgs("77001234567").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery("SELECT .+ FROM users WHERE phone").
		WithArgs("77001234567").
		WillReturnRows(sqlmock.NewRows([]string{"id", "phone", "name", "role", "is_active", "created_at"}).
			AddRow("user-id", "77001234567", "Test", "passenger", false, time.Now()))

	h := NewAuthHandler(authSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/verify", bytes.NewReader([]byte(`{"phone":"77001234567","code":"1234"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.VerifyOTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestAuthHandler_VerifyOTP_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	authSvc := service.NewAuthService(repository.NewUserRepository(db), repository.NewOTPRepository(db), &mockMessenger{})

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

	h := NewAuthHandler(authSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/verify", bytes.NewReader([]byte(`{"phone":"77001234567","code":"1234"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.VerifyOTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

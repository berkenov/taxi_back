package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestOTPRepository_StoreAndVerify(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewOTPRepository(db)
	phone := "77001234567"
	code := "1234"

	mock.ExpectExec("INSERT INTO otp_codes").
		WithArgs(phone, code, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Store(phone, code)
	if err != nil {
		t.Errorf("Store: %v", err)
	}

	expiresAt := time.Now().Add(10 * time.Minute)
	rows := sqlmock.NewRows([]string{"code", "expires_at"}).AddRow(code, expiresAt)
	mock.ExpectQuery("SELECT .+ FROM otp_codes").
		WithArgs(phone).
		WillReturnRows(rows)

	valid, err := repo.Verify(phone, code)
	if err != nil {
		t.Errorf("Verify: %v", err)
	}
	if !valid {
		t.Error("Verify: expected true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

func TestOTPRepository_Verify_WrongCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewOTPRepository(db)
	expiresAt := time.Now().Add(10 * time.Minute)
	rows := sqlmock.NewRows([]string{"code", "expires_at"}).AddRow("9999", expiresAt)

	mock.ExpectQuery("SELECT .+ FROM otp_codes").
		WithArgs("77001234567").
		WillReturnRows(rows)

	valid, err := repo.Verify("77001234567", "1234")
	if err != nil {
		t.Errorf("Verify: %v", err)
	}
	if valid {
		t.Error("Verify: expected false for wrong code")
	}
}

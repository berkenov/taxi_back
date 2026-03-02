package service

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"taxi/internal/models"
	"taxi/internal/repository"
)

func setupUserRepoMock(t *testing.T) (*repository.UserRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return repository.NewUserRepository(db), mock
}

func TestUserService_Register_Success(t *testing.T) {
	userRepo, mock := setupUserRepoMock(t)
	svc := NewUserService(userRepo)

	mock.ExpectQuery("SELECT .+ FROM users WHERE phone").
		WithArgs("77001234567").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO users").
		WithArgs(sqlmock.AnyArg(), "77001234567", "Иван", "passenger", false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	user, err := svc.Register("77001234567", "Иван", models.RolePassenger)
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if user.Phone != "77001234567" {
		t.Errorf("phone = %q", user.Phone)
	}
	if user.Role != models.RolePassenger {
		t.Errorf("role = %q", user.Role)
	}
}

func TestUserService_Register_PhoneExists(t *testing.T) {
	userRepo, mock := setupUserRepoMock(t)
	svc := NewUserService(userRepo)

	rows := sqlmock.NewRows([]string{"id", "phone", "name", "role", "is_active", "created_at"}).
		AddRow("id", "77001234567", "Old", "passenger", false, time.Now())

	mock.ExpectQuery("SELECT .+ FROM users WHERE phone").
		WithArgs("77001234567").
		WillReturnRows(rows)

	_, err := svc.Register("77001234567", "Иван", models.RolePassenger)
	if err != ErrPhoneExists {
		t.Errorf("Register: want ErrPhoneExists, got %v", err)
	}
}

func TestUserService_Register_InvalidPhone(t *testing.T) {
	svc := NewUserService(nil)

	_, err := svc.Register("123", "Иван", models.RolePassenger)
	if err != ErrInvalidPhone {
		t.Errorf("Register: want ErrInvalidPhone, got %v", err)
	}
}

func TestUserService_Register_InvalidRole(t *testing.T) {
	svc := NewUserService(nil)

	_, err := svc.Register("77001234567", "Иван", "admin")
	if err != ErrInvalidRole {
		t.Errorf("Register: want ErrInvalidRole, got %v", err)
	}
}

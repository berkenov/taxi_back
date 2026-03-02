package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"taxi/internal/models"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	user := &models.User{
		Phone:    "77001234567",
		Name:     "Test",
		Role:     models.RolePassenger,
		IsActive: false,
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(sqlmock.AnyArg(), "77001234567", "Test", "passenger", false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(user)
	if err != nil {
		t.Errorf("Create: %v", err)
	}
	if user.ID == "" {
		t.Error("Create: ID should be set")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	id := "550e8400-e29b-41d4-a716-446655440000"
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "phone", "name", "role", "is_active", "created_at"}).
		AddRow(id, "77001234567", "Test", "passenger", false, createdAt)

	mock.ExpectQuery("SELECT .+ FROM users WHERE id").
		WithArgs(id).
		WillReturnRows(rows)

	user, err := repo.GetByID(id)
	if err != nil {
		t.Errorf("GetByID: %v", err)
	}
	if user == nil {
		t.Fatal("GetByID: expected user")
	}
	if user.Phone != "77001234567" {
		t.Errorf("GetByID: phone = %q", user.Phone)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	mock.ExpectQuery("SELECT .+ FROM users WHERE id").
		WithArgs("00000000-0000-0000-0000-000000000000").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByID("00000000-0000-0000-0000-000000000000")
	if err != nil {
		t.Errorf("GetByID: %v", err)
	}
	if user != nil {
		t.Error("GetByID: expected nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

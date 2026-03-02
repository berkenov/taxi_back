package repository

import (
	"database/sql"
	"taxi/internal/models"

	"github.com/google/uuid"
)

// UserRepository handles user persistence
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user
func (r *UserRepository) Create(user *models.User) error {
	user.ID = uuid.New().String()
	query := `INSERT INTO users (id, phone, name, role, is_active) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, user.ID, user.Phone, user.Name, user.Role, user.IsActive)
	return err
}

// GetByID finds user by ID
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `SELECT id, phone, name, role, is_active, created_at FROM users WHERE id = $1`
	var u models.User
	err := r.db.QueryRow(query, id).Scan(
		&u.ID, &u.Phone, &u.Name, &u.Role, &u.IsActive, &u.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByPhone finds user by phone number
func (r *UserRepository) GetByPhone(phone string) (*models.User, error) {
	query := `SELECT id, phone, name, role, is_active, created_at FROM users WHERE phone = $1`
	var u models.User
	err := r.db.QueryRow(query, phone).Scan(
		&u.ID, &u.Phone, &u.Name, &u.Role, &u.IsActive, &u.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Update updates user by ID
func (r *UserRepository) Update(user *models.User) error {
	query := `UPDATE users SET phone = $1, name = $2, role = $3, is_active = $4 WHERE id = $5`
	result, err := r.db.Exec(query, user.Phone, user.Name, user.Role, user.IsActive, user.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Delete removes user by ID
func (r *UserRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

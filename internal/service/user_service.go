package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"taxi/internal/models"
	"taxi/internal/repository"
	"taxi/pkg/validator"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrPhoneExists     = errors.New("phone already registered")
	ErrInvalidRole     = errors.New("role must be passenger or driver")
	ErrInvalidPhone    = errors.New("invalid phone format")
	ErrNameRequired    = errors.New("name is required")
)

// UserService contains business logic for users
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// ValidatePhone checks if phone number is valid (format 7XXXXXXXXXX)
func (s *UserService) ValidatePhone(phone string) bool {
	return validator.ValidatePhone(phone)
}

// GetUserByPhone returns user by phone or nil if not found
func (s *UserService) GetUserByPhone(phone string) (*models.User, error) {
	normalized := validator.NormalizePhone(phone)
	return s.userRepo.GetByPhone(normalized)
}

// Register creates a new user (passenger or driver)
func (s *UserService) Register(phone, name, role string) (*models.User, error) {
	if !validator.ValidatePhone(phone) {
		return nil, ErrInvalidPhone
	}
	if name == "" {
		return nil, ErrNameRequired
	}
	if role != models.RolePassenger && role != models.RoleDriver {
		return nil, ErrInvalidRole
	}

	normalized := validator.NormalizePhone(phone)
	existing, err := s.userRepo.GetByPhone(normalized)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrPhoneExists
	}

	user := &models.User{
		Phone:    normalized,
		Name:     name,
		Role:     role,
		IsActive: false,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

// GetByID returns user by ID
func (s *UserService) GetByID(id string) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// Update updates user by ID
func (s *UserService) Update(id string, phone, name, role *string, isActive *bool) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	if phone != nil {
		if !validator.ValidatePhone(*phone) {
			return nil, ErrInvalidPhone
		}
		user.Phone = validator.NormalizePhone(*phone)
	}
	if name != nil {
		if *name == "" {
			return nil, ErrNameRequired
		}
		user.Name = *name
	}
	if role != nil {
		if *role != models.RolePassenger && *role != models.RoleDriver {
			return nil, ErrInvalidRole
		}
		user.Role = *role
	}
	if isActive != nil {
		user.IsActive = *isActive
	}

	if err := s.userRepo.Update(user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, ErrPhoneExists
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	return user, nil
}

// Delete removes user by ID
func (s *UserService) Delete(id string) error {
	err := s.userRepo.Delete(id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}
	return err
}

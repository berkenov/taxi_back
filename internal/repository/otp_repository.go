package repository

import (
	"database/sql"
	"time"
)

const OTPExpiryMinutes = 5

// OTPRepository handles OTP code storage
type OTPRepository struct {
	db *sql.DB
}

// NewOTPRepository creates a new OTPRepository
func NewOTPRepository(db *sql.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

// Store saves OTP for phone with expiry
func (r *OTPRepository) Store(phone, code string) error {
	expiresAt := time.Now().Add(OTPExpiryMinutes * time.Minute)
	query := `INSERT INTO otp_codes (phone, code, expires_at) VALUES ($1, $2, $3)
		ON CONFLICT (phone) DO UPDATE SET code = $2, expires_at = $3`
	_, err := r.db.Exec(query, phone, code, expiresAt)
	return err
}

// Verify checks if code matches and returns true if valid (not expired)
func (r *OTPRepository) Verify(phone, code string) (bool, error) {
	var storedCode string
	var expiresAt time.Time
	query := `SELECT code, expires_at FROM otp_codes WHERE phone = $1`
	err := r.db.QueryRow(query, phone).Scan(&storedCode, &expiresAt)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if time.Now().After(expiresAt) {
		return false, nil
	}
	return storedCode == code, nil
}

// Delete removes OTP after successful verification
func (r *OTPRepository) Delete(phone string) error {
	_, err := r.db.Exec(`DELETE FROM otp_codes WHERE phone = $1`, phone)
	return err
}

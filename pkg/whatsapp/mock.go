package whatsapp

import "log"

// MockMessenger logs messages instead of sending them (for tests)
type MockMessenger struct{}

// SendOTP logs the OTP instead of sending via API
func (m *MockMessenger) SendOTP(phone, code string) error {
	log.Printf("[MockMessenger] Would send OTP to %s: code=%s", phone, code)
	return nil
}

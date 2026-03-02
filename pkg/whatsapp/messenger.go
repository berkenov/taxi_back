package whatsapp

// Messenger defines the interface for sending messages (OTP, notifications).
// Allows mocking in tests without sending real messages.
type Messenger interface {
	SendOTP(phone, code string) error
}

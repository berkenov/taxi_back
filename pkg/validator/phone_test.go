package validator

import "testing"

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  bool
	}{
		{"valid 7XXXXXXXXXX", "77001234567", true},
		{"valid with plus", "+77001234567", true},
		{"valid 8 prefix", "87001234567", true},
		{"invalid too short", "7700123456", false},
		{"invalid too long", "770012345678", false},
		{"invalid letters", "7700abc4567", false},
		{"invalid empty", "", false},
		{"valid with spaces", " 77001234567 ", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePhone(tt.phone); got != tt.want {
				t.Errorf("ValidatePhone(%q) = %v, want %v", tt.phone, got, tt.want)
			}
		})
	}
}

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		phone string
		want  string
	}{
		{"77001234567", "77001234567"},
		{"+77001234567", "77001234567"},
		{"8 700 123 45 67", "77001234567"},
		{"7001234567", "77001234567"},
		{"", ""},
		{"abc", ""},
	}
	for _, tt := range tests {
		t.Run(tt.phone, func(t *testing.T) {
			if got := NormalizePhone(tt.phone); got != tt.want {
				t.Errorf("NormalizePhone(%q) = %q, want %q", tt.phone, got, tt.want)
			}
		})
	}
}

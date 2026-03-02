package validator

import (
	"regexp"
	"strings"
)

// PhoneRegex validates Kazakh phone format: 7XXXXXXXXXX (11 digits, starts with 7)
var PhoneRegex = regexp.MustCompile(`^7[0-9]{10}$`)

// ValidatePhone checks if the phone number is valid (format 7XXXXXXXXXX)
func ValidatePhone(phone string) bool {
	normalized := NormalizePhone(phone)
	return len(normalized) == 11 && PhoneRegex.MatchString(normalized)
}

// NormalizePhone converts phone to standard format 7XXXXXXXXXX (11 digits)
func NormalizePhone(phone string) string {
	digits := regexp.MustCompile(`[0-9]`).FindAllString(phone, -1)
	result := strings.Join(digits, "")
	if len(result) == 0 {
		return ""
	}
	// 87001234567 -> 77001234567 (8 is country code, replace with 7)
	if len(result) == 11 && result[0] == '8' {
		result = "7" + result[1:]
	}
	// +77001234567 or 77001234567 -> already correct
	if len(result) == 11 && result[0] == '7' {
		return result
	}
	// 7001234567 (10 digits, 70X format) -> 77001234567
	// 7700123456 (10 digits, 77X) is incomplete - return as is, will fail validation
	if len(result) == 10 && result[0] == '7' {
		if len(result) > 1 && result[1] == '0' {
			return "7" + result
		}
		return result
	}
	// 8012345678 (10 digits, 8+local) -> 77001234567 (8->7, add 7)
	if len(result) == 10 && result[0] == '8' {
		return "7" + result[1:]
	}
	// 0123456789 (10 digits) -> 7012345678? No - 7+10 digits needed
	if len(result) == 10 {
		return "7" + result
	}
	if len(result) == 11 {
		return result
	}
	return result
}

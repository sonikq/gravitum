package validator

import (
	"regexp"
	"strings"
)

// ValidEmail - validates email.
func ValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

// ValidGender - validating gender. (F/M/O)
func ValidGender(gender string) bool {
	if strings.ToUpper(gender) == "M" ||
		strings.ToUpper(gender) == "F" ||
		strings.ToUpper(gender) == "O" {
		return true
	}

	return false
}

// ValidAge - validating age. (1<=age<=150)
func ValidAge(age uint8) bool {
	if age <= 150 && age >= 1 {
		return true
	}

	return false
}

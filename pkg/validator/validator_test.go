package validator

import "testing"

// Tests
func TestValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"invalid_email@", false},
		{"@example.com", false},
		{"test@.com", false},
		{"test@domain.co", true},
		{"test@domain.toolongtld", false},
	}

	for _, test := range tests {
		result := ValidEmail(test.email)
		if result != test.valid {
			t.Errorf("Expected ValidEmail(%q) = %v, got %v", test.email, test.valid, result)
		}
	}
}

func TestValidGender(t *testing.T) {
	tests := []struct {
		gender string
		valid  bool
	}{
		{"M", true},
		{"F", true},
		{"O", true},
		{"m", true},
		{"f", true},
		{"o", true},
		{"X", false},
		{"", false},
	}

	for _, test := range tests {
		result := ValidGender(test.gender)
		if result != test.valid {
			t.Errorf("Expected ValidGender(%q) = %v, got %v", test.gender, test.valid, result)
		}
	}
}

func TestValidAge(t *testing.T) {
	tests := []struct {
		age   uint8
		valid bool
	}{
		{25, true},
		{1, true},
		{150, true},
		{0, false},
		{151, false},
		{100, true},
	}

	for _, test := range tests {
		result := ValidAge(test.age)
		if result != test.valid {
			t.Errorf("Expected ValidAge(%d) = %v, got %v", test.age, test.valid, result)
		}
	}
}

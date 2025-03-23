package models

import (
	"github.com/sonikq/gravitum_test_task/pkg/validator"
	"time"
)

type UserInfo struct {
	ID         int64      `json:"id"`
	Username   string     `json:"username"`
	FirstName  string     `json:"first_name"`
	MiddleName string     `json:"middle_name,omitempty"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email"`
	Gender     string     `json:"gender"`
	Age        uint8      `json:"age"`
	EndDate    *time.Time `json:"-"`
}

func (uf *UserInfo) Validate() error {
	if !validator.ValidEmail(uf.Email) {
		return ErrInvalidEmail
	}

	if !validator.ValidGender(uf.Gender) {
		return ErrInvalidGender
	}

	if !validator.ValidAge(uf.Age) {
		return ErrInvalidAge
	}

	return nil
}

package password

import (
	"errors"
	"unicode"
)

var (
	ErrTooShort    = errors.New("password must be at least 6 characters")
	ErrNoUppercase = errors.New("password must contain at least one uppercase letter")
	ErrNoLowercase = errors.New("password must contain at least one lowercase letter")
	ErrNoDigit     = errors.New("password must contain at least one digit")
	ErrNoSpecial   = errors.New("password must contain at least one special character")
)

func Validate(password string) error {
	if len(password) < 6 {
		return ErrTooShort
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrNoUppercase
	}
	if !hasLower {
		return ErrNoLowercase
	}
	if !hasDigit {
		return ErrNoDigit
	}
	if !hasSpecial {
		return ErrNoSpecial
	}
	return nil
}

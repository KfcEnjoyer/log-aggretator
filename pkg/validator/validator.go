package validator

import (
	"errors"
	"strings"
)

type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) CheckError(err error) {
	if err != nil {
		v.Errors[err.Error()] = err.Error()
	}
}

func (v *Validator) ValidateUsername(username string) (bool, error) {
	if username == "" {
		err := errors.New("username is required")
		return false, err
	} else if len(username) < 3 {
		err := errors.New("username cannot be less than 3 characters")
		return false, err
	}

	return true, nil
}

func (v *Validator) ValidatePassword(p string) (bool, error) {
	if p == "" {
		err := errors.New("password cannot be empty")
		return false, err
	} else if len(p) < 6 {
		err := errors.New("password cannot be less than 6 characters")
		return false, err
	} else if !strings.ContainsAny(p, "_@!%") {
		err := errors.New("password should contain a special character e.g. _@!%")
		return false, err
	}

	return true, nil
}

package validator

import "errors"

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

func (v *Validator) ValidateParams(username, password string) (bool, error) {
	if username == "" {
		err := errors.New("username is required")
		return false, err
	} else if password == "" {
		err := errors.New("password is required")
		return false, err
	}

	return true, nil
}

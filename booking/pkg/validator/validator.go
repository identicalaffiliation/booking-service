package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{v: validator.New()}
}

func (v *Validator) Validate(input any) error {
	return v.v.Struct(input)
}

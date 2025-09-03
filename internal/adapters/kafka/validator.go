package kafka

import (
	"github.com/go-playground/validator/v10"
)

// Validator оборачивает библиотеку валидации
type Validator struct {
	validate *validator.Validate
}

// NewValidator создает Validator
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Validate выполняет валидацию над структурой
func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

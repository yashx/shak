package shak

import (
	"github.com/yashx/shak/types"
	"github.com/yashx/shak/validation"
)

// Satisfies reports whether value passes all rules.
func Satisfies[T any](value T, rules ...types.Rule[T]) bool {
	err := Validate(value, rules...)
	return err == nil
}

// Validate runs rules against value and returns the first error, or nil.
func Validate[T any](value T, rules ...types.Rule[T]) *types.ValidationError {
	v := validation.Value("", value, rules...)
	return v.Validate()
}

// SatisfiesValidation reports whether v produces no errors.
func SatisfiesValidation[T validation.Validatable](v T) bool {
	err := v.Validation().Validate()
	return err == nil
}

// RunValidation executes v and returns the first error, or nil.
func RunValidation[T validation.Validatable](v T) *types.ValidationError {
	return v.Validation().Validate()
}

// ValidateAll runs all rules against value and returns every error.
func ValidateAll[T any](value T, rules ...types.Rule[T]) []*types.ValidationError {
	v := validation.Value("", value, rules...)
	return v.ValidateAll()
}

// RunFullValidation executes v and returns every error.
func RunFullValidation[T validation.Validatable](v T) []*types.ValidationError {
	return v.Validation().ValidateAll()
}

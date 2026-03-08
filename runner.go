package shak

import (
	"github.com/yashx/shak/types"
	"github.com/yashx/shak/validation"
)

func Satisfies[T any](value T, rules ...types.Rule[T]) bool {
	err := Validate(value, rules...)
	return err == nil
}

func Validate[T any](value T, rules ...types.Rule[T]) *types.ValidationError {
	v := validation.Value(".", value, rules...)
	return v.Validate()
}

func RunValidation(v validation.Validation) *types.ValidationError {
	return v.Validate()
}

func ValidateAll[T any](value T, rules ...types.Rule[T]) []*types.ValidationError {
	v := validation.Value(".", value, rules...)
	return v.ValidateAll()
}

func RunFullValidation(v validation.Validation) []*types.ValidationError {
	return v.ValidateAll()
}

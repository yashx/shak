package validation

import "github.com/yashx/shak/types"

// Validatable is implemented by types that describe their own validation rules.
// Call Validation() once per validation cycle to get a consistent snapshot.
type Validatable interface {
	Validation() Validation
}

// Validation is an executable validation produced by Value, NewValidations, or Nested.
// Validate stops at the first error; ValidateAll collects every error.
type Validation interface {
	Validate() *types.ValidationError
	ValidateAll() []*types.ValidationError
}

type valueValidation[T any] struct {
	name  string
	value T
	rules []types.Rule[T]
}

// Value creates a Validation that runs rules against a single named value.
func Value[T any](name string, value T, rules ...types.Rule[T]) Validation {
	return &valueValidation[T]{
		name:  name,
		value: value,
		rules: rules,
	}
}

func (v *valueValidation[T]) Validate() *types.ValidationError {
	for _, r := range v.rules {
		if r == nil {
			continue
		}
		if err := r(v.value); err != nil {
			return types.NewValidationErrorWithParentName(v.name, err)
		}
	}
	return nil
}

func (v *valueValidation[T]) ValidateAll() []*types.ValidationError {
	var errors []*types.ValidationError
	for _, r := range v.rules {
		if r == nil {
			continue
		}
		if err := r(v.value); err != nil {
			errors = append(errors, types.NewValidationErrorWithParentName(v.name, err))
		}
	}
	return errors
}

type validationsHolder struct {
	validations []Validation
}

// NewValidations groups multiple Validation entries, running them in order.
func NewValidations(validations ...Validation) Validation {
	return &validationsHolder{
		validations: validations,
	}
}

func (s *validationsHolder) Validate() *types.ValidationError {
	for _, v := range s.validations {
		if v == nil {
			continue
		}
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (s *validationsHolder) ValidateAll() []*types.ValidationError {
	var errors []*types.ValidationError
	for _, v := range s.validations {
		if v == nil {
			continue
		}
		ers := v.ValidateAll()
		errors = append(errors, ers...)
	}
	return errors
}

type nestedValidation[T Validatable] struct {
	name  string
	value T
	rules []types.Rule[T]
}

// Nested creates a Validation for a named field whose type implements Validatable.
// Outer rules run first; if they pass, the field's own Validation() runs next.
// All error paths are prefixed with name.
func Nested[T Validatable](name string, value T, rules ...types.Rule[T]) Validation {
	return &nestedValidation[T]{
		name:  name,
		value: value,
		rules: rules,
	}
}

func (v *nestedValidation[T]) Validate() *types.ValidationError {
	for _, r := range v.rules {
		if r == nil {
			continue
		}
		if err := r(v.value); err != nil {
			return types.NewValidationErrorWithParentName(v.name, err)
		}
	}
	err := v.value.Validation().Validate()
	if err != nil {
		return types.NewValidationErrorWithParentName(v.name, err)
	}
	return nil
}

func (v *nestedValidation[T]) ValidateAll() []*types.ValidationError {
	var errors []*types.ValidationError
	for _, r := range v.rules {
		if r == nil {
			continue
		}
		if err := r(v.value); err != nil {
			errors = append(errors, types.NewValidationErrorWithParentName(v.name, err))
		}
	}

	ers := v.value.Validation().ValidateAll()
	for _, e := range ers {
		errors = append(errors, types.NewValidationErrorWithParentName(v.name, e))
	}

	return errors
}

package validation

import "github.com/yashx/shak/types"

// Validatable is implemented by types that know how to describe their own
// validation. The Validation method returns a Validation that captures the
// current field values and their rules — call it once per validation cycle
// to get a consistent snapshot.
//
// Implement this on your domain structs so they can be passed to
// RunValidation / RunFullValidation or composed into parent validations via
// NewValidations.
type Validatable interface {
	Validation() Validation
}

// Validation represents a named, executable validation — the result of
// describing what to check and against which rules. It is produced by
// Value or NewValidations and then executed by calling Validate or
// ValidateAll.
//
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

// NewValidations groups multiple Validation entries under a named parent.
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

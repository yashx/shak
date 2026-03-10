package rule

import (
	"github.com/yashx/shak/types"
)

// IsZeroValue validates that the value equals its type's zero value.
func IsZeroValue[T comparable]() types.Rule[T] {
	return func(value T) *types.ValidationError {
		var zero T
		if value == zero {
			return nil
		}
		return types.NewValidationError("%v is not the zero value", value)
	}
}

// IsNotZeroValue validates that the value does not equal its type's zero value.
func IsNotZeroValue[T comparable]() types.Rule[T] {
	return func(value T) *types.ValidationError {
		var zero T
		if value != zero {
			return nil
		}
		return types.NewValidationError("value must not be the zero value")
	}
}

// Equal validates that the value equals eq.
func Equal[T comparable](eq T) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if value == eq {
			return nil
		}
		return types.NewValidationError("%v is not equal to %v", value, eq)
	}
}

// NotEqual validates that the value does not equal ne.
func NotEqual[T comparable](ne T) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if value != ne {
			return nil
		}
		return types.NewValidationError("%v must not equal %v", value, ne)
	}
}

// EqualToAny validates that the value equals at least one of the provided values.
func EqualToAny[T comparable](eq ...T) types.Rule[T] {
	return func(value T) *types.ValidationError {
		for _, v := range eq {
			if value == v {
				return nil
			}
		}
		return types.NewValidationError("%v does not match any of the allowed values %v", value, eq)
	}
}

// NotEqualToAny validates that the value does not equal any of the provided values.
func NotEqualToAny[T comparable](ne ...T) types.Rule[T] {
	return func(value T) *types.ValidationError {
		for _, v := range ne {
			if value == v {
				return types.NewValidationError("%v matches a forbidden value in %v", value, ne)
			}
		}
		return nil
	}
}

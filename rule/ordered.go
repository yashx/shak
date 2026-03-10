package rule

import (
	"cmp"

	"github.com/yashx/shak/types"
)

// Positive validates that the value is greater than zero.
func Positive[T cmp.Ordered]() types.Rule[T] {
	return func(value T) *types.ValidationError {
		var zero T
		if value > zero {
			return nil
		}
		return types.NewValidationError("value %v is not positive", value)
	}
}

// PositiveOrZero validates that the value is greater than or equal to zero.
func PositiveOrZero[T cmp.Ordered]() types.Rule[T] {
	return func(value T) *types.ValidationError {
		var zero T
		if value >= zero {
			return nil
		}
		return types.NewValidationError("value %v is not positive or zero", value)
	}
}

// Negative validates that the value is less than zero.
func Negative[T cmp.Ordered]() types.Rule[T] {
	return func(value T) *types.ValidationError {
		var zero T
		if value < zero {
			return nil
		}
		return types.NewValidationError("value %v is not negative", value)
	}
}

// NegativeOrZero validates that the value is less than or equal to zero.
func NegativeOrZero[T cmp.Ordered]() types.Rule[T] {
	return func(value T) *types.ValidationError {
		var zero T
		if value <= zero {
			return nil
		}
		return types.NewValidationError("value %v is not negative or zero", value)
	}
}

// Min validates that the value is greater than or equal to min.
func Min[T cmp.Ordered](min T) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if min <= value {
			return nil
		}
		return types.NewValidationError("value %v is less than minimum %v", value, min)
	}
}

// Max validates that the value is less than or equal to max.
func Max[T cmp.Ordered](max T) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if max >= value {
			return nil
		}
		return types.NewValidationError("value %v is greater than maximum %v", value, max)
	}
}

// Range validates that the value is between min and max inclusive.
func Range[T cmp.Ordered](min, max T) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if min <= value && value <= max {
			return nil
		}
		return types.NewValidationError("value %v is not between %v and %v", value, min, max)
	}
}

// NotInRange validates that the value is not between min and max inclusive.
func NotInRange[T cmp.Ordered](min, max T) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if value < min || value > max {
			return nil
		}
		return types.NewValidationError("value %v must not be between %v and %v", value, min, max)
	}
}

// LessThan validates that the value is strictly less than le.
func LessThan[T cmp.Ordered](le T) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if value < le {
			return nil
		}
		return types.NewValidationError("value %v is not less than %v", value, le)
	}
}

// LessThanEqual validates that the value is less than or equal to le.
func LessThanEqual[T cmp.Ordered](le T) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if value <= le {
			return nil
		}
		return types.NewValidationError("value %v is not less than or equal to %v", value, le)
	}
}

// GreaterThan validates that the value is strictly greater than ge.
func GreaterThan[T cmp.Ordered](ge T) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if value > ge {
			return nil
		}
		return types.NewValidationError("value %v is not greater than %v", value, ge)
	}
}

// GreaterThanEqual validates that the value is greater than or equal to ge.
func GreaterThanEqual[T cmp.Ordered](ge T) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if value >= ge {
			return nil
		}
		return types.NewValidationError("value %v is not greater than or equal to %v", value, ge)
	}
}

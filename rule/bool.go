package rule

import (
	"github.com/yashx/shak/types"
)

// IsTrue validates that the value is true.
func IsTrue[T ~bool]() types.Rule[T] {
	return func(b T) *types.ValidationError {
		if b {
			return nil
		}
		return types.NewValidationError("value is not true")
	}
}

// IsFalse validates that the value is false.
func IsFalse[T ~bool]() types.Rule[T] {
	return func(b T) *types.ValidationError {
		if !b {
			return nil
		}
		return types.NewValidationError("value is not false")
	}
}

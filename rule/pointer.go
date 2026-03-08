package rule

import "github.com/yashx/shak/types"

// IsNil validates that the pointer is nil.
func IsNil[T any]() types.Rule[*T] {
	return func(p *T) *types.ValidationError {
		if p == nil {
			return nil
		}
		return types.NewValidationError("value is not nil")
	}
}

// IsNotNil validates that the pointer is not nil.
func IsNotNil[T any]() types.Rule[*T] {
	return func(p *T) *types.ValidationError {
		if p != nil {
			return nil
		}
		return types.NewValidationError("value is nil")
	}
}

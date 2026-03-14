package rule

import (
	"github.com/yashx/shak/types"
	"github.com/yashx/shak/validation"
)

// IsValid validates that the value passes its own Validation() rules.
func IsValid[T validation.Validatable]() types.Rule[T] {
	return func(value T) *types.ValidationError {
		return value.Validation().Validate()
	}
}

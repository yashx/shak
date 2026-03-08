package rule

import (
	"fmt"
	"regexp"

	"github.com/yashx/shak/types"
)

// ForEach validates each slice element against the provided rule.
func ForEach[S ~[]T, T any](s *S, rule types.Rule[T]) types.Rule[S] {
	// s is only passed to help compiler with type inference
	return func(value S) *types.ValidationError {
		for i, v := range value {
			if err := rule(v); err != nil {
				return types.NewValidationErrorWithPath(fmt.Sprintf("[%v]", i), err.Message)
			}
		}
		return nil
	}
}

// LengthSlice validates that the slice has exactly l elements.
func LengthSlice[S ~[]T, T any](s *S, l int) types.Rule[S] {
	return func(value S) *types.ValidationError {
		if len(value) == l {
			return nil
		}
		return types.NewValidationError("length is %d, expected %d", len(value), l)
	}
}

// Match validates that the byte slice matches the provided regular expression.
func Match[T ~[]byte](re *regexp.Regexp) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if re.Match(value) {
			return nil
		}
		return types.NewValidationError("value does not match the required format")
	}
}

package rule

import (
	"fmt"

	"github.com/yashx/shak/types"
)

// ForEachValue validates each map value against the provided rule.
func ForEachValue[M ~map[K]T, K comparable, T any](m *M, rule types.Rule[T]) types.Rule[M] {
	// m is only passed to help compiler with type inference
	return func(value M) *types.ValidationError {
		for k, v := range value {
			if err := rule(v); err != nil {
				return types.NewValidationErrorWithPath(fmt.Sprintf("[%v]", k), err.Message)
			}
		}
		return nil
	}
}

// ForKey validates the value at a specific map key against the provided rule.
// Returns an error if the key is not present in the map.
func ForKey[M ~map[K]T, K comparable, T any](m *M, key K, rule types.Rule[T]) types.Rule[M] {
	// m is only passed to help compiler with type inference
	return func(m M) *types.ValidationError {
		value, ok := m[key]
		if !ok {
			return types.NewValidationErrorWithPath(fmt.Sprintf("[%v]", key), "value is missing")
		}
		if err := rule(value); err != nil {
			return types.NewValidationErrorWithPath(fmt.Sprintf("[%v]", key), err.Message)
		}
		return nil
	}
}

// LengthMap validates that the map has exactly l entries.
func LengthMap[M ~map[K]T, K comparable, T any](m *M, l int) types.Rule[M] {
	return func(value M) *types.ValidationError {
		if len(value) == l {
			return nil
		}
		return types.NewValidationError("length is %d, expected %d", len(value), l)
	}
}

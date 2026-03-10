package rule

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"

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

// ContainsElem validates that the slice contains v.
func ContainsElem[S ~[]T, T comparable](s *S, v T) types.Rule[S] {
	return func(s S) *types.ValidationError {
		if slices.Contains(s, v) {
			return nil
		}
		return types.NewValidationError("%v not found", v)
	}
}

// EqualSlice validates that the slice is equal to s2.
func EqualSlice[S ~[]T, T comparable](s2 S) types.Rule[S] {
	return func(value S) *types.ValidationError {
		if slices.Equal(value, s2) {
			return nil
		}
		return types.NewValidationError("not equal to %v", s2)
	}
}

// EqualSliceFunc validates that the slice is equal to s2 using a custom equality function.
// s is only passed to help the compiler with type inference.
func EqualSliceFunc[S1 ~[]T1, S2 ~[]T2, T1, T2 any](s *S1, s2 S2, eq func(T1, T2) bool) types.Rule[S1] {
	return func(value S1) *types.ValidationError {
		if slices.EqualFunc(value, s2, eq) {
			return nil
		}
		return types.NewValidationError("not equal to %v", s2)
	}
}

// IsSorted validates that the slice is sorted in ascending order.
// s is only passed to help the compiler with type inference.
func IsSorted[S ~[]T, T cmp.Ordered](s *S) types.Rule[S] {
	return func(value S) *types.ValidationError {
		if slices.IsSorted(value) {
			return nil
		}
		return types.NewValidationError("slice is not sorted")
	}
}

// IsSortedFunc validates that the slice is sorted according to the provided comparison function.
// s is only passed to help the compiler with type inference.
func IsSortedFunc[S ~[]T, T any](s *S, cmpFn func(T, T) int) types.Rule[S] {
	return func(value S) *types.ValidationError {
		if slices.IsSortedFunc(value, cmpFn) {
			return nil
		}
		return types.NewValidationError("slice is not sorted")
	}
}

// MatchBytes validates that the byte slice matches the provided regular expression.
func MatchBytes[T ~[]byte](re *regexp.Regexp) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if re.Match(value) {
			return nil
		}
		return types.NewValidationError("value does not match the required format")
	}
}

// NotMatchBytes validates that the byte slice does not match the provided regular expression.
func NotMatchBytes[T ~[]byte](re *regexp.Regexp) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if !re.Match(value) {
			return nil
		}
		return types.NewValidationError("value must not match the pattern %q", re.String())
	}
}

// MinLengthSlice validates that the slice has at least l elements.
func MinLengthSlice[S ~[]T, T any](s *S, l int) types.Rule[S] {
	return func(value S) *types.ValidationError {
		if len(value) >= l {
			return nil
		}
		return types.NewValidationError("length is %d, expected at least %d", len(value), l)
	}
}

// MaxLengthSlice validates that the slice has at most l elements.
func MaxLengthSlice[S ~[]T, T any](s *S, l int) types.Rule[S] {
	return func(value S) *types.ValidationError {
		if len(value) <= l {
			return nil
		}
		return types.NewValidationError("length is %d, expected at most %d", len(value), l)
	}
}

// NotContainsElem validates that the slice does not contain v.
func NotContainsElem[S ~[]T, T comparable](s *S, v T) types.Rule[S] {
	return func(value S) *types.ValidationError {
		if !slices.Contains(value, v) {
			return nil
		}
		return types.NewValidationError("%v must not appear in slice", v)
	}
}

// Unique validates that all elements in the slice are distinct.
func Unique[S ~[]T, T comparable](s *S) types.Rule[S] {
	return func(value S) *types.ValidationError {
		seen := make(map[T]struct{}, len(value))
		for _, v := range value {
			if _, exists := seen[v]; exists {
				return types.NewValidationError("duplicate value %v found", v)
			}
			seen[v] = struct{}{}
		}
		return nil
	}
}

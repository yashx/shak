package rule

import (
	"github.com/yashx/shak/types"
)

// Or combines multiple rules into one, returning nil if any rule passes.
// Returns the last error if all rules fail. Nil rules are skipped.
func Or[T any](rules ...types.Rule[T]) types.Rule[T] {
	return func(value T) *types.ValidationError {
		var last *types.ValidationError
		for _, rule := range rules {
			if rule == nil {
				continue
			}
			if err := rule(value); err == nil {
				return nil
			} else {
				last = err
			}
		}
		return last
	}
}

// And combines multiple rules into one, returning the first error encountered.
func And[T any](rules ...types.Rule[T]) types.Rule[T] {
	return func(value T) *types.ValidationError {
		for _, rule := range rules {
			if rule == nil {
				continue
			}
			if err := rule(value); err != nil {
				return err
			}
		}
		return nil
	}
}

// Ref lifts value rules into a pointer rule.
func Ref[T any](rules ...types.Rule[T]) types.Rule[*T] {
	return func(value *T) *types.ValidationError {
		if value == nil {
			return types.NewValidationError("value is nil")
		}
		for _, rule := range rules {
			if rule == nil {
				continue
			}
			if err := rule(*value); err != nil {
				return err
			}
		}
		return nil
	}
}

// Deref lifts pointer rules into a value rule by taking the address of the value.
func Deref[T any](rules ...types.Rule[*T]) types.Rule[T] {
	return func(value T) *types.ValidationError {
		for _, rule := range rules {
			if rule == nil {
				continue
			}
			if err := rule(&value); err != nil {
				return err
			}
		}
		return nil
	}
}

// If returns the rule when condition is true, or nil when false.
// Nil rules are safely skipped by And, Else, and the validation layer — do not call the result directly.
func If[T any](condition bool, conditionalRule types.Rule[T]) types.Rule[T] {
	if !condition {
		return nil
	}
	return conditionalRule
}

// Package shaktest provides test helpers that validate values with shak.
// Import it only from _test.go files.
package shaktest

import (
	"testing"

	"github.com/yashx/shak"
	"github.com/yashx/shak/types"
	"github.com/yashx/shak/validation"
)

// Require runs all rules against value, reports every error and, if any, stops the test.
// Like t.FailNow, it must be called from the goroutine running the test.
func Require[T any](t testing.TB, value T, rules ...types.Rule[T]) {
	t.Helper()
	if !Assert(t, value, rules...) {
		t.FailNow()
	}
}

// Assert runs all rules against value, reports every error and reports whether value passed.
// The test continues either way.
func Assert[T any](t testing.TB, value T, rules ...types.Rule[T]) bool {
	t.Helper()
	return report(t, shak.ValidateAll(value, rules...))
}

// RequireValidation executes v, reports every error and, if any, stops the test.
// Like t.FailNow, it must be called from the goroutine running the test.
func RequireValidation(t testing.TB, v validation.Validatable) {
	t.Helper()
	if !AssertValidation(t, v) {
		t.FailNow()
	}
}

// AssertValidation executes v, reports every error and reports whether v passed.
// The test continues either way.
func AssertValidation(t testing.TB, v validation.Validatable) bool {
	t.Helper()
	return report(t, shak.RunFullValidation(v))
}

func report(t testing.TB, errs []*types.ValidationError) bool {
	t.Helper()
	for _, err := range errs {
		t.Error(err)
	}
	return len(errs) == 0
}

package shaktest_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/yashx/shak/rule"
	"github.com/yashx/shak/shaktest"
	"github.com/yashx/shak/types"
	"github.com/yashx/shak/validation"
)

// fakeTB records reported errors and FailNow calls instead of failing the test.
type fakeTB struct {
	testing.TB
	errors          []string
	failNows        int
	errorsAtFailNow int
}

func (f *fakeTB) Helper() {}

func (f *fakeTB) Error(args ...any) {
	f.errors = append(f.errors, fmt.Sprint(args...))
}

func (f *fakeTB) FailNow() {
	f.failNows++
	f.errorsAtFailNow = len(f.errors)
}

type address struct {
	Street string
	City   string
}

func (a address) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("street", a.Street, rule.NotBlank[string]()),
		validation.Value("city", a.City, rule.NotBlank[string]()),
	)
}

type order struct {
	ID      int
	Address address
}

func (o order) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("ID", o.ID, rule.Positive[int]()),
		validation.Nested("address", o.Address),
	)
}

func checkErrors(t *testing.T, f *fakeTB, want []string) {
	t.Helper()
	if !slices.Equal(f.errors, want) {
		t.Errorf("errors: got %q, want %q", f.errors, want)
	}
}

func TestAssertPass(t *testing.T) {
	f := &fakeTB{TB: t}
	if !shaktest.Assert(f, 5, rule.Positive[int](), rule.Max(10)) {
		t.Error("Assert returned false for a valid value")
	}
	checkErrors(t, f, nil)
	if f.failNows != 0 {
		t.Errorf("FailNow called %d times, want 0", f.failNows)
	}
}

func TestAssertReportsEveryError(t *testing.T) {
	f := &fakeTB{TB: t}
	if shaktest.Assert(f, -5, rule.Positive[int](), rule.Min(0), rule.Max(10)) {
		t.Error("Assert returned true for an invalid value")
	}
	checkErrors(t, f, []string{
		"value -5 is not positive",
		"value -5 is less than minimum 0",
	})
	if f.failNows != 0 {
		t.Errorf("FailNow called %d times, want 0", f.failNows)
	}
}

func TestAssertSkipsNilRules(t *testing.T) {
	f := &fakeTB{TB: t}
	if !shaktest.Assert(f, 5, nil, rule.Positive[int](), rule.If(false, rule.Max(0))) {
		t.Error("Assert returned false for a valid value")
	}
	checkErrors(t, f, nil)
}

func TestAssertNoRules(t *testing.T) {
	f := &fakeTB{TB: t}
	if !shaktest.Assert[int](f, 5) {
		t.Error("Assert returned false with no rules")
	}
	checkErrors(t, f, nil)
}

func TestRequirePass(t *testing.T) {
	f := &fakeTB{TB: t}
	shaktest.Require(f, "abc", rule.NotBlank[string]())
	checkErrors(t, f, nil)
	if f.failNows != 0 {
		t.Errorf("FailNow called %d times, want 0", f.failNows)
	}
}

func TestRequireReportsEveryErrorThenStops(t *testing.T) {
	f := &fakeTB{TB: t}
	shaktest.Require(f, "", rule.NotBlank[string](), rule.MinLengthString[string](3))
	checkErrors(t, f, []string{
		"value must not be blank",
		"length is 0 bytes, expected at least 3",
	})
	if f.failNows != 1 {
		t.Errorf("FailNow called %d times, want 1", f.failNows)
	}
	if f.errorsAtFailNow != 2 {
		t.Errorf("FailNow called after %d errors, want 2", f.errorsAtFailNow)
	}
}

func TestAssertValidationPass(t *testing.T) {
	f := &fakeTB{TB: t}
	o := order{ID: 1, Address: address{Street: "1 Oak Ave", City: "Austin"}}
	if !shaktest.AssertValidation(f, o) {
		t.Error("AssertValidation returned false for a valid value")
	}
	checkErrors(t, f, nil)
}

func TestAssertValidationReportsNestedPaths(t *testing.T) {
	f := &fakeTB{TB: t}
	o := order{ID: 0, Address: address{Street: "", City: ""}}
	if shaktest.AssertValidation(f, o) {
		t.Error("AssertValidation returned true for an invalid value")
	}
	checkErrors(t, f, []string{
		"ID: value 0 is not positive",
		"address.street: value must not be blank",
		"address.city: value must not be blank",
	})
	if f.failNows != 0 {
		t.Errorf("FailNow called %d times, want 0", f.failNows)
	}
}

func TestRequireValidationPass(t *testing.T) {
	f := &fakeTB{TB: t}
	shaktest.RequireValidation(f, address{Street: "1 Oak Ave", City: "Austin"})
	checkErrors(t, f, nil)
	if f.failNows != 0 {
		t.Errorf("FailNow called %d times, want 0", f.failNows)
	}
}

func TestRequireValidationReportsEveryErrorThenStops(t *testing.T) {
	f := &fakeTB{TB: t}
	shaktest.RequireValidation(f, address{})
	checkErrors(t, f, []string{
		"street: value must not be blank",
		"city: value must not be blank",
	})
	if f.failNows != 1 {
		t.Errorf("FailNow called %d times, want 1", f.failNows)
	}
	if f.errorsAtFailNow != 2 {
		t.Errorf("FailNow called after %d errors, want 2", f.errorsAtFailNow)
	}
}

func TestCustomRuleError(t *testing.T) {
	f := &fakeTB{TB: t}
	even := func(n int) *types.ValidationError {
		if n%2 == 0 {
			return nil
		}
		return types.NewValidationError("%d is odd", n)
	}
	shaktest.Assert(f, 3, even)
	checkErrors(t, f, []string{"3 is odd"})
}

// TestWithRealT runs the helpers against a real *testing.T with passing values.
func TestWithRealT(t *testing.T) {
	shaktest.Require(t, 5, rule.Positive[int]())
	shaktest.Assert(t, "abc", rule.NotBlank[string]())
	shaktest.RequireValidation(t, address{Street: "1 Oak Ave", City: "Austin"})
	shaktest.AssertValidation(t, order{ID: 1, Address: address{Street: "1 Oak Ave", City: "Austin"}})
}

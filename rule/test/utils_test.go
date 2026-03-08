package test

import (
	"testing"

	"github.com/yashx/shak/rule"
)

// And

func TestMergeIntPassesAll(t *testing.T) {
	r := rule.And(rule.Min[int](1), rule.Max[int](100))
	if err := r(50); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMergeIntFailsFirst(t *testing.T) {
	// 0 fails Min(1), never reaches Max
	r := rule.And(rule.Min[int](1), rule.Max[int](100))
	if err := r(0); err == nil {
		t.Errorf("Error expected when first rule fails")
	}
}

func TestMergeIntFailsSecond(t *testing.T) {
	// 200 passes Min(1), fails Max(100)
	r := rule.And(rule.Min[int](1), rule.Max[int](100))
	if err := r(200); err == nil {
		t.Errorf("Error expected when second rule fails")
	}
}

func TestMergeStopsAtFirstError(t *testing.T) {
	// Both rules fail for -5; only the first error should be returned
	r := rule.And(rule.Min[int](1), rule.Max[int](100))
	err := r(-5)
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	// Min(1) fires first, so message should mention minimum
	r2 := rule.Min[int](1)
	expected := r2(-5)
	if err.Message != expected.Message {
		t.Errorf("Expected first rule's error, got: %v", err.Message)
	}
}

func TestMergeThreeRulesFloat64(t *testing.T) {
	// positive, at least 1.0, at most 9.9
	r := rule.And(rule.Positive[float64](), rule.Min[float64](1.0), rule.Max[float64](9.9))
	if err := r(5.5); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMergeThreeRulesFloat64Error(t *testing.T) {
	r := rule.And(rule.Positive[float64](), rule.Min[float64](1.0), rule.Max[float64](9.9))
	if err := r(-1.0); err == nil {
		t.Errorf("Error expected for negative float64")
	}
}

func TestMergeSkipsNilRule(t *testing.T) {
	// If(false, ...) returns nil — And should skip it without panicking
	r := rule.And(rule.If(false, rule.Max[int](0)), rule.Min[int](1))
	if err := r(5); err != nil {
		t.Errorf("No errors expected when nil rule is skipped, got: %v", err)
	}
}

func TestMergeAllNilRules(t *testing.T) {
	r := rule.And(rule.If[int](false, rule.Max[int](0)), rule.If[int](false, rule.Min[int](10)))
	if err := r(50); err != nil {
		t.Errorf("No errors expected when all rules are nil, got: %v", err)
	}
}

func TestMergeCustomStringType(t *testing.T) {
	type Tag string
	r := rule.And(rule.IsNotZeroValue[Tag](), rule.EqualIgnoreCase[Tag]("go"))
	if err := r(Tag("GO")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMergeCustomStringTypeError(t *testing.T) {
	type Tag string
	r := rule.And(rule.IsNotZeroValue[Tag](), rule.EqualIgnoreCase[Tag]("go"))
	if err := r(Tag("")); err == nil {
		t.Errorf("Error expected for empty tag (zero value)")
	}
}

// And — zero arguments

func TestMergeNoRules(t *testing.T) {
	r := rule.And[int]()
	if err := r(42); err != nil {
		t.Errorf("No errors expected with no rules, got: %v", err)
	}
}

// WithMessage

func TestWithMessageOverridesErrorMessage(t *testing.T) {
	r := rule.Positive[int]().WithMessage("must be a positive number")
	err := r(-1)
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Message != "must be a positive number" {
		t.Errorf("Expected overridden message, got: %v", err.Message)
	}
}

func TestWithMessagePassesWhenRulePasses(t *testing.T) {
	r := rule.Positive[int]().WithMessage("must be a positive number")
	if err := r(5); err != nil {
		t.Errorf("No errors expected when rule passes, got: %v", err)
	}
}

func TestWithMessageFloat64(t *testing.T) {
	r := rule.Max[float64](100.0).WithMessage("value exceeds the maximum allowed")
	err := r(200.0)
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Message != "value exceeds the maximum allowed" {
		t.Errorf("Expected overridden message, got: %v", err.Message)
	}
}

func TestWithMessageCustomType(t *testing.T) {
	type Username string
	r := rule.IsNotZeroValue[Username]().WithMessage("username is required")
	err := r(Username(""))
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Message != "username is required" {
		t.Errorf("Expected overridden message, got: %v", err.Message)
	}
}

func TestWithMessageFormattedMessage(t *testing.T) {
	limit := 10
	r := rule.Max[int](limit).WithMessage("value must not exceed %d", limit)
	err := r(20)
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Message != "value must not exceed 10" {
		t.Errorf("Expected formatted message, got: %v", err.Message)
	}
}

// Else

func TestElseNonNilRuleIsUsed(t *testing.T) {
	// non-nil rule.Else(other) returns the original rule, not other
	r := rule.Min[int](5).Else(rule.Max[int](0))
	// Min(5) is active — value 10 should pass
	if err := r(10); err != nil {
		t.Errorf("No errors expected, original rule should be active, got: %v", err)
	}
}

func TestElseNonNilRuleIgnoresElseBranch(t *testing.T) {
	// original rule fails — Else is not a fallback for failures, only for nil rules
	r := rule.Min[int](5).Else(rule.Max[int](100))
	if err := r(3); err == nil {
		t.Errorf("Error expected from original rule, Else does not act as fallback")
	}
}

func TestElseNilRuleUsesElseBranch(t *testing.T) {
	// If(false, ...) produces nil; .Else(rule) should activate the else rule
	r := rule.If(false, rule.Min[int](5)).Else(rule.Positive[int]())
	if err := r(5); err != nil {
		t.Errorf("No errors expected, else branch should be active, got: %v", err)
	}
}

func TestElseNilRuleElseBranchFails(t *testing.T) {
	r := rule.If(false, rule.Min[int](5)).Else(rule.Positive[int]())
	if err := r(-1); err == nil {
		t.Errorf("Error expected from else branch rule")
	}
}

func TestElseWithString(t *testing.T) {
	r := rule.If(false, rule.Equal[string]("admin")).Else(rule.IsNotZeroValue[string]())
	if err := r("hello"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// If

func TestIfTrueRulePasses(t *testing.T) {
	r := rule.If(true, rule.Positive[int]())
	if err := r(5); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIfTrueRuleFails(t *testing.T) {
	r := rule.If(true, rule.Positive[int]())
	if err := r(-1); err == nil {
		t.Errorf("Error expected when condition is true and value fails the rule")
	}
}

func TestIfFalseReturnsNil(t *testing.T) {
	// If(false, ...) must return nil so callers like And and Else can detect it
	r := rule.If(false, rule.Min[int](10))
	if r != nil {
		t.Errorf("Expected nil rule when condition is false")
	}
}

func TestIfTrueReturnsRule(t *testing.T) {
	inner := rule.Min[int](10)
	r := rule.If(true, inner)
	if r == nil {
		t.Errorf("Expected non-nil rule when condition is true")
	}
}

func TestIfElseTrueBranch(t *testing.T) {
	// condition true → Min(10) applies; value 15 passes
	v := 15
	r := rule.If(v > 0, rule.Min[int](10)).Else(rule.Negative[int]())
	if err := r(v); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIfElseTrueBranchFails(t *testing.T) {
	// condition true → Min(10) applies; value 5 fails
	v := 5
	r := rule.If(v > 0, rule.Min[int](10)).Else(rule.Negative[int]())
	if err := r(v); err == nil {
		t.Errorf("Error expected when true-branch rule fails")
	}
}

func TestIfElseFalseBranch(t *testing.T) {
	// condition false → Negative applies; value -3 passes
	v := -3
	r := rule.If(v > 0, rule.Min[int](10)).Else(rule.Negative[int]())
	if err := r(v); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIfElseFalseBranchFails(t *testing.T) {
	// condition false → Negative applies; value 5 fails Negative
	v := -1
	r := rule.If(v > 0, rule.Min[int](10)).Else(rule.Negative[int]())
	if err := r(5); err == nil {
		t.Errorf("Error expected when false-branch rule fails")
	}
}

func TestIfWithStringType(t *testing.T) {
	isAdmin := true
	r := rule.If(isAdmin, rule.Equal[string]("admin"))
	if err := r("admin"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIfWithStringTypeFalseViaMerge(t *testing.T) {
	// condition false — skipped via And; "guest" would fail Equal("admin") if applied
	isAdmin := false
	r := rule.And(rule.If(isAdmin, rule.Equal[string]("admin")))
	if err := r("guest"); err != nil {
		t.Errorf("No errors expected when condition is false, got: %v", err)
	}
}

func TestIfWithCustomType(t *testing.T) {
	type Score int
	aboveThreshold := true
	r := rule.If(aboveThreshold, rule.Min[Score](50))
	if err := r(Score(75)); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIfWithCustomTypeError(t *testing.T) {
	type Score int
	aboveThreshold := true
	r := rule.If(aboveThreshold, rule.Min[Score](50))
	if err := r(Score(30)); err == nil {
		t.Errorf("Error expected when condition is true and value fails rule")
	}
}

// Or

func TestOrFirstRulePasses(t *testing.T) {
	// Min(1) passes for 5 — Or should return nil without checking Max(3)
	r := rule.Or(rule.Min[int](1), rule.Max[int](3))
	if err := r(5); err != nil {
		t.Errorf("No errors expected when first rule passes, got: %v", err)
	}
}

func TestOrSecondRulePasses(t *testing.T) {
	// Min(10) fails for 5, Max(6) passes — Or should return nil
	r := rule.Or(rule.Min[int](10), rule.Max[int](6))
	if err := r(5); err != nil {
		t.Errorf("No errors expected when second rule passes, got: %v", err)
	}
}

func TestOrAllRulesFail(t *testing.T) {
	// 5 fails Min(10) and fails Max(3)
	r := rule.Or(rule.Min[int](10), rule.Max[int](3))
	if err := r(5); err == nil {
		t.Errorf("Error expected when all rules fail")
	}
}

func TestOrReturnsLastErrorWhenAllFail(t *testing.T) {
	last := rule.Max[int](3)
	r := rule.Or(rule.Min[int](10), last)
	err := r(5)
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	expected := last(5)
	if err.Message != expected.Message {
		t.Errorf("Expected last rule's error message, got: %v", err.Message)
	}
}

func TestOrSkipsNilRules(t *testing.T) {
	// nil rule is skipped; Positive passes for 5
	r := rule.Or(rule.If(false, rule.Max[int](0)), rule.Positive[int]())
	if err := r(5); err != nil {
		t.Errorf("No errors expected when non-nil rule passes, got: %v", err)
	}
}

func TestOrAllNilRulesPass(t *testing.T) {
	r := rule.Or(rule.If[int](false, rule.Max[int](0)), rule.If[int](false, rule.Min[int](10)))
	if err := r(50); err != nil {
		t.Errorf("No errors expected when all rules are nil, got: %v", err)
	}
}

func TestOrNoRules(t *testing.T) {
	r := rule.Or[int]()
	if err := r(42); err != nil {
		t.Errorf("No errors expected with no rules, got: %v", err)
	}
}

func TestOrFloat64(t *testing.T) {
	// passes if negative OR greater than 100
	r := rule.Or(rule.Negative[float64](), rule.Min[float64](100))
	if err := r(-1.0); err != nil {
		t.Errorf("No errors expected for negative value, got: %v", err)
	}
	if err := r(150.0); err != nil {
		t.Errorf("No errors expected for value above 100, got: %v", err)
	}
}

func TestOrFloat64AllFail(t *testing.T) {
	r := rule.Or(rule.Negative[float64](), rule.Min[float64](100))
	if err := r(50.0); err == nil {
		t.Errorf("Error expected when value is neither negative nor above 100")
	}
}

func TestOrCustomType(t *testing.T) {
	type Status string
	r := rule.Or(rule.Equal[Status]("active"), rule.Equal[Status]("pending"))
	if err := r(Status("active")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
	if err := r(Status("pending")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestOrCustomTypeAllFail(t *testing.T) {
	type Status string
	r := rule.Or(rule.Equal[Status]("active"), rule.Equal[Status]("pending"))
	if err := r(Status("deleted")); err == nil {
		t.Errorf("Error expected for status not in allowed set")
	}
}

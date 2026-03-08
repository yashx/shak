package shak

import (
	"testing"

	"github.com/yashx/shak/rule"
	"github.com/yashx/shak/validation"
)

// Validate

func TestValidateNoError(t *testing.T) {
	err := Validate(10, rule.Min(5), rule.Max(15))
	if err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestValidateReturnsFirstError(t *testing.T) {
	err := Validate(100, rule.Min(5), rule.Max(15))
	if err == nil {
		t.Errorf("Error expected, got nil")
	}
}

func TestValidateNoRules(t *testing.T) {
	err := Validate(42)
	if err != nil {
		t.Errorf("No errors expected with no rules, got: %v", err)
	}
}

// Satisfies

func TestSatisfiesTrue(t *testing.T) {
	if !Satisfies(10, rule.Min(5), rule.Max(15)) {
		t.Errorf("Expected Satisfies to return true")
	}
}

func TestSatisfiesFalse(t *testing.T) {
	if Satisfies(100, rule.Max(15)) {
		t.Errorf("Expected Satisfies to return false")
	}
}

func TestSatisfiesNoRules(t *testing.T) {
	if !Satisfies(42) {
		t.Errorf("Expected Satisfies to return true with no rules")
	}
}

// ValidateAll

func TestValidateAllNoErrors(t *testing.T) {
	errs := ValidateAll(10, rule.Min(5), rule.Max(15))
	if len(errs) != 0 {
		t.Errorf("No errors expected, got: %v", errs)
	}
}

func TestValidateAllCollectsAllErrors(t *testing.T) {
	// both rules fail: 100 > Max(15) and 100 < Min(200)
	errs := ValidateAll(100, rule.Min(200), rule.Max(15))
	if len(errs) != 2 {
		t.Errorf("Expected 2 errors, got: %d", len(errs))
	}
}

func TestValidateAllVsValidateDifference(t *testing.T) {
	// Validate stops at first error; ValidateAll collects all
	err := Validate(100, rule.Min(200), rule.Max(15))
	errs := ValidateAll(100, rule.Min(200), rule.Max(15))
	if err == nil {
		t.Errorf("Validate: error expected")
	}
	if len(errs) <= 1 {
		t.Errorf("ValidateAll: expected more than 1 error, got: %d", len(errs))
	}
}

// RunValidation / RunFullValidation

type sample struct {
	X int
	Y int
}

func (s sample) Validation() validation.Validation {
	return validation.NewValidations(
		"sample",
		validation.Value("X", s.X, rule.Max(10)),
		validation.Value("Y", s.Y, rule.Min(10)),
	)
}

func TestRunValidationNoError(t *testing.T) {
	s := sample{X: 5, Y: 20}
	err := RunValidation(s.Validation())
	if err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestRunValidationReturnsFirstError(t *testing.T) {
	// X fails Max(10) and Y fails Min(10); only first error returned
	s := sample{X: 11, Y: 1}
	err := RunValidation(s.Validation())
	if err == nil {
		t.Errorf("Error expected, got nil")
	}
}

func TestRunFullValidationNoErrors(t *testing.T) {
	s := sample{X: 5, Y: 20}
	errs := RunFullValidation(s.Validation())
	for _, e := range errs {
		if e != nil {
			t.Errorf("No errors expected, got: %v", e)
		}
	}
}

func TestRunFullValidationCollectsAllErrors(t *testing.T) {
	// both X and Y fail
	s := sample{X: 11, Y: 1}
	errs := RunFullValidation(s.Validation())
	var count int
	for _, e := range errs {
		if e != nil {
			count++
		}
	}
	if count != 2 {
		t.Errorf("Expected 2 errors, got: %d", count)
	}
}

func TestRunFullValidationVsRunValidationDifference(t *testing.T) {
	s := sample{X: 11, Y: 1}
	err := RunValidation(s.Validation())
	errs := RunFullValidation(s.Validation())
	if err == nil {
		t.Errorf("RunValidation: error expected")
	}
	var count int
	for _, e := range errs {
		if e != nil {
			count++
		}
	}
	if count <= 1 {
		t.Errorf("RunFullValidation: expected more than 1 error, got: %d", count)
	}
}

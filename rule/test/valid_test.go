package test

import (
	"testing"

	"github.com/yashx/shak/rule"
	"github.com/yashx/shak/types"
	"github.com/yashx/shak/validation"
)

// stubValidatable controls pass/fail via a pre-set error, for simple IsValid tests.
type stubValidatable struct {
	err *types.ValidationError
}

func (s stubValidatable) Validation() validation.Validation {
	return validation.Value("stub", 0, func(_ int) *types.ValidationError {
		return s.err
	})
}

// product is a real-world-style Validatable with actual field rules.
type product struct {
	Name  string
	Price float64
}

func (p product) Validation() validation.Validation {
	return validation.NewValidations("product",
		validation.Value("name", p.Name, rule.IsNotZeroValue[string]()),
		validation.Value("price", p.Price, rule.Positive[float64]()),
	)
}

// IsValid

func TestIsValidPasses(t *testing.T) {
	r := rule.IsValid[stubValidatable]()
	if err := r(stubValidatable{err: nil}); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIsValidFails(t *testing.T) {
	r := rule.IsValid[stubValidatable]()
	if err := r(stubValidatable{err: types.NewValidationError("inner rule failed")}); err == nil {
		t.Errorf("Error expected when inner Validate returns an error")
	}
}

func TestIsValidPropagatesErrorMessage(t *testing.T) {
	r := rule.IsValid[stubValidatable]()
	err := r(stubValidatable{err: types.NewValidationError("something is wrong")})
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Message != "something is wrong" {
		t.Errorf("Expected original error message, got: %v", err.Message)
	}
}

func TestIsValidRealStructPasses(t *testing.T) {
	r := rule.IsValid[product]()
	if err := r(product{Name: "Widget", Price: 9.99}); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIsValidRealStructFailsEmptyName(t *testing.T) {
	r := rule.IsValid[product]()
	if err := r(product{Name: "", Price: 9.99}); err == nil {
		t.Errorf("Error expected for empty name")
	}
}

func TestIsValidRealStructFailsNegativePrice(t *testing.T) {
	r := rule.IsValid[product]()
	if err := r(product{Name: "Widget", Price: -1.0}); err == nil {
		t.Errorf("Error expected for negative price")
	}
}

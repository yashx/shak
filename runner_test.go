package shak

import (
	"testing"

	"github.com/yashx/shak/rule"
	"github.com/yashx/shak/types"
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

// SatisfiesValidation

func TestSatisfiesValidationTrue(t *testing.T) {
	s := sample{X: 5, Y: 20}
	if !SatisfiesValidation(s.Validation()) {
		t.Errorf("Expected SatisfiesValidation to return true")
	}
}

func TestSatisfiesValidationFalse(t *testing.T) {
	s := sample{X: 11, Y: 1}
	if SatisfiesValidation(s.Validation()) {
		t.Errorf("Expected SatisfiesValidation to return false")
	}
}

func TestSatisfiesValidationStopsAtFirstError(t *testing.T) {
	// Both fields fail, but SatisfiesValidation just returns false — not which field failed
	s := sample{X: 100, Y: 0}
	if SatisfiesValidation(s.Validation()) {
		t.Errorf("Expected SatisfiesValidation to return false when any field fails")
	}
}

// Nested validation

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

type person struct {
	Name    string
	Age     int
	Address address
}

func (p person) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("name", p.Name, rule.NotBlank[string]()),
		validation.Value("age", p.Age, rule.Positive[int](), rule.Range(0, 120)),
		validation.Nested("address", p.Address),
	)
}

func TestRunValidationNestedNoError(t *testing.T) {
	p := person{Name: "Alice", Age: 30, Address: address{Street: "1 Main St", City: "Springfield"}}
	if err := RunValidation(p.Validation()); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestRunValidationNestedInnerFails(t *testing.T) {
	p := person{Name: "Alice", Age: 30, Address: address{Street: "", City: "Springfield"}}
	if err := RunValidation(p.Validation()); err == nil {
		t.Errorf("Error expected for blank street")
	}
}

func TestRunValidationNestedOuterFails(t *testing.T) {
	p := person{Name: "", Age: 30, Address: address{Street: "1 Main St", City: "Springfield"}}
	if err := RunValidation(p.Validation()); err == nil {
		t.Errorf("Error expected for blank name")
	}
}

func TestRunFullValidationNestedCollectsAllErrors(t *testing.T) {
	// name blank + address both fields blank = 3 errors total
	p := person{Name: "", Age: 30, Address: address{Street: "", City: ""}}
	errs := RunFullValidation(p.Validation())
	if len(errs) != 3 {
		t.Errorf("Expected 3 errors, got: %d", len(errs))
	}
}

func TestRunFullValidationNestedOuterAndInnerFail(t *testing.T) {
	// age out of range + street blank = 2 errors across two levels
	p := person{Name: "Alice", Age: 200, Address: address{Street: "", City: "Springfield"}}
	errs := RunFullValidation(p.Validation())
	if len(errs) != 2 {
		t.Errorf("Expected 2 errors, got: %d", len(errs))
	}
}

func TestSatisfiesValidationNested(t *testing.T) {
	p := person{Name: "Alice", Age: 30, Address: address{Street: "1 Main St", City: "Springfield"}}
	if !SatisfiesValidation(p.Validation()) {
		t.Errorf("Expected SatisfiesValidation to return true")
	}
}

func TestSatisfiesValidationNestedFails(t *testing.T) {
	p := person{Name: "Alice", Age: 30, Address: address{Street: "", City: "Springfield"}}
	if SatisfiesValidation(p.Validation()) {
		t.Errorf("Expected SatisfiesValidation to return false for invalid nested struct")
	}
}

// Custom rules on top of nested validation
//
// Nested runs your custom outer rules first, then the struct's own Validation().
// This lets you layer business rules on top of structural validity.

var allowedCities = []string{"Austin", "Denver", "Portland"}

func notServicedCity(a address) *types.ValidationError {
	for _, c := range allowedCities {
		if a.City == c {
			return nil
		}
	}
	return types.NewValidationError("delivery not available in %q", a.City)
}

type delivery struct {
	From        string
	Destination address
}

func (d delivery) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("from", d.From, rule.NotBlank[string]()),
		// notServicedCity is a custom outer rule; address.Validation() also runs
		validation.Nested("destination", d.Destination, notServicedCity),
	)
}

func TestCustomRuleOnNestedValid(t *testing.T) {
	d := delivery{From: "warehouse", Destination: address{Street: "1 Oak Ave", City: "Austin"}}
	if err := RunValidation(d.Validation()); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestCustomRuleOnNestedCityNotServiced(t *testing.T) {
	// custom outer rule fails; address itself is structurally valid
	d := delivery{From: "warehouse", Destination: address{Street: "1 Oak Ave", City: "Miami"}}
	if err := RunValidation(d.Validation()); err == nil {
		t.Errorf("Error expected: city not in serviced list")
	}
}

func TestCustomRuleOnNestedBothFail(t *testing.T) {
	// custom outer rule fails AND inner struct is invalid (blank street)
	d := delivery{From: "warehouse", Destination: address{Street: "", City: "Miami"}}
	errs := RunFullValidation(d.Validation())
	if len(errs) != 2 {
		t.Errorf("Expected 2 errors (city + street), got: %d", len(errs))
	}
}

// Value vs Nested: using Value on a Validatable type does NOT run its Validation()
//
// If you use validation.Value instead of validation.Nested, only the rules you
// explicitly pass are checked — the struct's own Validation() is never called.
// This lets you define a completely different set of rules for the same type
// depending on context.

func cityNotBlank(a address) *types.ValidationError {
	if a.City != "" {
		return nil
	}
	return types.NewValidationError("city must not be blank")
}

type shipment struct {
	TrackingID string
	Drop       address
}

func (s shipment) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("trackingID", s.TrackingID, rule.NotBlank[string]()),
		// Value with a custom rule — only cityNotBlank runs, not address.Validation()
		// so a blank street is never checked
		validation.Value("drop", s.Drop, cityNotBlank),
	)
}

func TestValueRunsOnlyExplicitRules(t *testing.T) {
	// Street is blank, which would fail address.Validation() — but we used Value,
	// so only cityNotBlank runs and the shipment is considered valid
	s := shipment{TrackingID: "TRK-001", Drop: address{Street: "", City: "Austin"}}
	if err := RunValidation(s.Validation()); err != nil {
		t.Errorf("No errors expected: blank street is not checked via Value, got: %v", err)
	}
}

func TestValueExplicitRuleFails(t *testing.T) {
	// The custom cityNotBlank rule fires; street being blank is still irrelevant
	s := shipment{TrackingID: "TRK-001", Drop: address{Street: "", City: ""}}
	if err := RunValidation(s.Validation()); err == nil {
		t.Errorf("Error expected: cityNotBlank rule should fail")
	}
}

package test

import (
	"testing"

	"github.com/yashx/shak/rule"
	"github.com/yashx/shak/types"
)

// --- Custom type and rules ---

type Person struct {
	Name string
	Age  int
}

// Custom rule on value receiver: Age must be non-negative and Name non-empty.
func validPerson() types.Rule[Person] {
	return func(p Person) *types.ValidationError {
		if p.Name == "" {
			return types.NewValidationError("name cannot be empty")
		}
		if p.Age < 0 {
			return types.NewValidationError("age %d cannot be negative", p.Age)
		}
		return nil
	}
}

// Custom rule on pointer receiver: Age must be at least 18.
func adultPersonPtr() types.Rule[*Person] {
	return func(p *Person) *types.ValidationError {
		if p.Age < 18 {
			return types.NewValidationError("age %d is below 18", p.Age)
		}
		return nil
	}
}

// --- Direct custom rule tests ---

func TestCustomRuleValid(t *testing.T) {
	r := validPerson()
	if err := r(Person{Name: "Alice", Age: 30}); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestCustomRuleEmptyName(t *testing.T) {
	r := validPerson()
	if err := r(Person{Name: "", Age: 25}); err == nil {
		t.Errorf("Error expected for empty name")
	}
}

func TestCustomRuleNegativeAge(t *testing.T) {
	r := validPerson()
	if err := r(Person{Name: "Bob", Age: -1}); err == nil {
		t.Errorf("Error expected for negative age")
	}
}

func TestCustomPtrRuleValid(t *testing.T) {
	r := adultPersonPtr()
	p := Person{Name: "Alice", Age: 20}
	if err := r(&p); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestCustomPtrRuleError(t *testing.T) {
	r := adultPersonPtr()
	p := Person{Name: "Alice", Age: 15}
	if err := r(&p); err == nil {
		t.Errorf("Error expected for age below 18")
	}
}

// --- Ref: value rule → pointer rule ---

func TestRefConvertsValueRuleToPointerRule(t *testing.T) {
	// validPerson is written for Person (value); Ref lifts it to *Person
	r := rule.Ref(validPerson())
	p := Person{Name: "Alice", Age: 25}
	if err := r(&p); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestRefPointerRuleFailsOnBadValue(t *testing.T) {
	r := rule.Ref(validPerson())
	p := Person{Name: "", Age: 25}
	if err := r(&p); err == nil {
		t.Errorf("Error expected when underlying value fails the rule")
	}
}

func TestRefMultipleValueRules(t *testing.T) {
	// Combine two value rules and lift them together to a pointer rule
	nameRule := func(p Person) *types.ValidationError {
		if p.Name == "" {
			return types.NewValidationError("name cannot be empty")
		}
		return nil
	}
	ageRule := func(p Person) *types.ValidationError {
		if p.Age < 0 {
			return types.NewValidationError("age cannot be negative")
		}
		return nil
	}
	r := rule.Ref(types.Rule[Person](nameRule), types.Rule[Person](ageRule))
	p := Person{Name: "Bob", Age: 30}
	if err := r(&p); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestRefMultipleValueRulesError(t *testing.T) {
	nameRule := func(p Person) *types.ValidationError {
		if p.Name == "" {
			return types.NewValidationError("name cannot be empty")
		}
		return nil
	}
	ageRule := func(p Person) *types.ValidationError {
		if p.Age < 0 {
			return types.NewValidationError("age cannot be negative")
		}
		return nil
	}
	r := rule.Ref(types.Rule[Person](nameRule), types.Rule[Person](ageRule))
	p := Person{Name: "Bob", Age: -5}
	if err := r(&p); err == nil {
		t.Errorf("Error expected when value fails one of the lifted rules")
	}
}

// --- Deref: pointer rule → value rule ---

func TestDerefConvertsPointerRuleToValueRule(t *testing.T) {
	// adultPersonPtr is written for *Person; Deref brings it down to Person
	r := rule.Deref(adultPersonPtr())
	if err := r(Person{Name: "Alice", Age: 20}); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestDerefValueRuleFailsOnBadValue(t *testing.T) {
	r := rule.Deref(adultPersonPtr())
	if err := r(Person{Name: "Alice", Age: 16}); err == nil {
		t.Errorf("Error expected when value fails the deref'd rule")
	}
}

func TestDerefMultiplePointerRules(t *testing.T) {
	nameRule := func(p *Person) *types.ValidationError {
		if p.Name == "" {
			return types.NewValidationError("name cannot be empty")
		}
		return nil
	}
	r := rule.Deref(types.Rule[*Person](nameRule), adultPersonPtr())
	if err := r(Person{Name: "Alice", Age: 20}); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestDerefMultiplePointerRulesError(t *testing.T) {
	nameRule := func(p *Person) *types.ValidationError {
		if p.Name == "" {
			return types.NewValidationError("name cannot be empty")
		}
		return nil
	}
	r := rule.Deref(types.Rule[*Person](nameRule), adultPersonPtr())
	if err := r(Person{Name: "Alice", Age: 15}); err == nil {
		t.Errorf("Error expected when value fails one of the deref'd rules")
	}
}

// --- Ref + Deref roundtrip ---

func TestRefThenDerefRoundtrip(t *testing.T) {
	// validPerson() is a value rule; Ref lifts it to *Person; Deref brings it back to Person.
	// The final rule should behave identically to the original.
	original := validPerson()
	roundtrip := rule.Deref(rule.Ref(original))

	valid := Person{Name: "Alice", Age: 25}
	if err := roundtrip(valid); err != nil {
		t.Errorf("No errors expected after roundtrip, got: %v", err)
	}
}

func TestRefThenDerefRoundtripError(t *testing.T) {
	original := validPerson()
	roundtrip := rule.Deref(rule.Ref(original))

	invalid := Person{Name: "", Age: 25}
	if err := roundtrip(invalid); err == nil {
		t.Errorf("Error expected after roundtrip for invalid value")
	}
}

func TestDerefThenRefRoundtrip(t *testing.T) {
	// adultPersonPtr() is a pointer rule; Deref brings it to Person; Ref lifts it back to *Person.
	original := adultPersonPtr()
	roundtrip := rule.Ref(rule.Deref(original))

	p := Person{Name: "Alice", Age: 20}
	if err := roundtrip(&p); err != nil {
		t.Errorf("No errors expected after roundtrip, got: %v", err)
	}
}

func TestDerefThenRefRoundtripError(t *testing.T) {
	original := adultPersonPtr()
	roundtrip := rule.Ref(rule.Deref(original))

	p := Person{Name: "Alice", Age: 15}
	if err := roundtrip(&p); err == nil {
		t.Errorf("Error expected after roundtrip for invalid value")
	}
}

package tests

import (
	"testing"

	"github.com/yashx/shak/rule"
	"github.com/yashx/shak/types"
	"github.com/yashx/shak/validation"
)

// --- Value / valueValidation ---

func TestValueName(t *testing.T) {
	v := validation.Value("age", 25, rule.Min[int](0))
	if v.Name() != "age" {
		t.Errorf("Expected name 'age', got: %v", v.Name())
	}
}

func TestValueValidateNoError(t *testing.T) {
	v := validation.Value("age", 25, rule.Min[int](0), rule.Max[int](100))
	if err := v.Validate(); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestValueValidateReturnsFirstError(t *testing.T) {
	v := validation.Value("age", -1, rule.Min[int](0), rule.Max[int](100))
	err := v.Validate()
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
}

func TestValueValidateErrorPath(t *testing.T) {
	v := validation.Value("age", -1, rule.Min[int](0))
	err := v.Validate()
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Path != "age" {
		t.Errorf("Expected path 'age', got: %v", err.Path)
	}
}

func TestValueValidateStopsAtFirstError(t *testing.T) {
	// both rules fail; only the first error should be returned
	firstRuleFired := false
	secondRuleFired := false
	r1 := types.Rule[int](func(v int) *types.ValidationError {
		firstRuleFired = true
		return types.NewValidationError("first rule failed")
	})
	r2 := types.Rule[int](func(v int) *types.ValidationError {
		secondRuleFired = true
		return types.NewValidationError("second rule failed")
	})
	v := validation.Value("field", 0, r1, r2)
	v.Validate()
	if !firstRuleFired {
		t.Errorf("Expected first rule to be called")
	}
	if secondRuleFired {
		t.Errorf("Expected second rule to be skipped after first error")
	}
}

func TestValueValidateSkipsNilRule(t *testing.T) {
	v := validation.Value("score", 10, nil, rule.Max[int](100))
	if err := v.Validate(); err != nil {
		t.Errorf("No errors expected when nil rule is skipped, got: %v", err)
	}
}

func TestValueValidateNoRules(t *testing.T) {
	v := validation.Value("x", 42)
	if err := v.Validate(); err != nil {
		t.Errorf("No errors expected with no rules, got: %v", err)
	}
}

func TestValueValidateAllNoErrors(t *testing.T) {
	v := validation.Value("age", 25, rule.Min[int](0), rule.Max[int](100))
	errs := v.ValidateAll()
	if len(errs) != 0 {
		t.Errorf("No errors expected, got: %v", errs)
	}
}

func TestValueValidateAllCollectsAllErrors(t *testing.T) {
	// value 200 fails both Max(100) and Min(201) — wait, let's use two independently failing rules
	r1 := types.Rule[int](func(v int) *types.ValidationError {
		return types.NewValidationError("first error")
	})
	r2 := types.Rule[int](func(v int) *types.ValidationError {
		return types.NewValidationError("second error")
	})
	v := validation.Value("field", 0, r1, r2)
	errs := v.ValidateAll()
	if len(errs) != 2 {
		t.Errorf("Expected 2 errors, got: %d", len(errs))
	}
}

func TestValueValidateAllErrorPaths(t *testing.T) {
	r1 := types.Rule[int](func(v int) *types.ValidationError {
		return types.NewValidationError("bad value")
	})
	v := validation.Value("count", 0, r1)
	errs := v.ValidateAll()
	if len(errs) != 1 {
		t.Fatalf("Expected 1 error, got: %d", len(errs))
	}
	if errs[0].Path != "count" {
		t.Errorf("Expected path 'count', got: %v", errs[0].Path)
	}
}

// --- NewValidations / validationsHolder ---

func TestNewValidationsName(t *testing.T) {
	v := validation.NewValidations("myStruct")
	if v.Name() != "myStruct" {
		t.Errorf("Expected name 'myStruct', got: %v", v.Name())
	}
}

func TestNewValidationsValidateNoError(t *testing.T) {
	type S struct{ X, Y int }
	s := S{X: 5, Y: 15}
	v := validation.NewValidations("s",
		validation.Value("X", s.X, rule.Max[int](10)),
		validation.Value("Y", s.Y, rule.Min[int](10)),
	)
	if err := v.Validate(); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNewValidationsValidateReturnsFirstError(t *testing.T) {
	type S struct{ X, Y int }
	s := S{X: 20, Y: 1}
	v := validation.NewValidations("s",
		validation.Value("X", s.X, rule.Max[int](10)),
		validation.Value("Y", s.Y, rule.Min[int](10)),
	)
	if err := v.Validate(); err == nil {
		t.Errorf("Error expected, got nil")
	}
}

func TestNewValidationsValidateNestedPath(t *testing.T) {
	// error path should be "user.age", not just "age"
	type User struct{ Age int }
	u := User{Age: -1}
	v := validation.NewValidations("user",
		validation.Value("age", u.Age, rule.Min[int](0)),
	)
	err := v.Validate()
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Path != "user.age" {
		t.Errorf("Expected path 'user.age', got: %v", err.Path)
	}
}

func TestNewValidationsValidateSkipsNilValidation(t *testing.T) {
	type S struct{ X int }
	s := S{X: 5}
	v := validation.NewValidations("s",
		nil,
		validation.Value("X", s.X, rule.Max[int](10)),
	)
	if err := v.Validate(); err != nil {
		t.Errorf("No errors expected when nil validation is skipped, got: %v", err)
	}
}

func TestNewValidationsValidateAllCollectsAllErrors(t *testing.T) {
	type S struct{ X, Y int }
	s := S{X: 20, Y: 1}
	v := validation.NewValidations("s",
		validation.Value("X", s.X, rule.Max[int](10)),
		validation.Value("Y", s.Y, rule.Min[int](10)),
	)
	errs := v.ValidateAll()
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

func TestNewValidationsValidateAllNestedPaths(t *testing.T) {
	type S struct{ X int }
	s := S{X: 20}
	v := validation.NewValidations("s",
		validation.Value("X", s.X, rule.Max[int](10)),
	)
	errs := v.ValidateAll()
	if len(errs) != 1 {
		t.Fatalf("Expected 1 error, got: %d", len(errs))
	}
	if errs[0].Path != "s.X" {
		t.Errorf("Expected path 's.X', got: %v", errs[0].Path)
	}
}

func TestNewValidationsValidateAllNoErrors(t *testing.T) {
	type S struct{ X int }
	s := S{X: 5}
	v := validation.NewValidations("s",
		validation.Value("X", s.X, rule.Max[int](10)),
	)
	errs := v.ValidateAll()
	if len(errs) != 0 {
		t.Errorf("No errors expected, got: %v", errs)
	}
}

// --- Validatable interface ---

type Address struct {
	Street string
	City   string
}

func (a Address) Validation() validation.Validation {
	return validation.NewValidations("address",
		validation.Value("street", a.Street, rule.IsNotZeroValue[string]()),
		validation.Value("city", a.City, rule.IsNotZeroValue[string]()),
	)
}

type Order struct {
	ID     int
	Amount float64
}

func (o Order) Validation() validation.Validation {
	return validation.NewValidations("order",
		validation.Value("ID", o.ID, rule.Positive[int]()),
		validation.Value("amount", o.Amount, rule.Positive[float64]()),
	)
}

func TestValidatableValidateNoError(t *testing.T) {
	a := Address{Street: "123 Main St", City: "Springfield"}
	if err := a.Validation().Validate(); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestValidatableValidateError(t *testing.T) {
	a := Address{Street: "", City: "Springfield"}
	if err := a.Validation().Validate(); err == nil {
		t.Errorf("Error expected for empty street")
	}
}

func TestValidatableValidateErrorPath(t *testing.T) {
	a := Address{Street: "", City: "Springfield"}
	err := a.Validation().Validate()
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Path != "address.street" {
		t.Errorf("Expected path 'address.street', got: %v", err.Path)
	}
}

func TestValidatableValidateAllCollectsErrors(t *testing.T) {
	// both fields empty
	a := Address{Street: "", City: ""}
	errs := a.Validation().ValidateAll()
	if len(errs) != 2 {
		t.Errorf("Expected 2 errors, got: %d", len(errs))
	}
}

func TestValidatableDifferentTypes(t *testing.T) {
	o := Order{ID: 1, Amount: 49.99}
	if err := o.Validation().Validate(); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestValidatableDifferentTypesError(t *testing.T) {
	o := Order{ID: -1, Amount: 49.99}
	err := o.Validation().Validate()
	if err == nil {
		t.Fatalf("Error expected for negative ID")
	}
	if err.Path != "order.ID" {
		t.Errorf("Expected path 'order.ID', got: %v", err.Path)
	}
}

func TestValidatableImplementsInterface(t *testing.T) {
	// compile-time check that both types satisfy the Validatable interface
	var _ validation.Validatable = Address{}
	var _ validation.Validatable = Order{}
}

// --- Nested ---

type LineItem struct {
	Name     string
	Quantity int
}

func (l LineItem) Validation() validation.Validation {
	return validation.NewValidations("lineItem",
		validation.Value("name", l.Name, rule.IsNotZeroValue[string]()),
		validation.Value("quantity", l.Quantity, rule.Positive[int]()),
	)
}

func TestNestedName(t *testing.T) {
	item := LineItem{Name: "widget", Quantity: 3}
	v := validation.Nested("item", item)
	if v.Name() != "item" {
		t.Errorf("Expected name 'item', got: %v", v.Name())
	}
}

func TestNestedValidateNoError(t *testing.T) {
	item := LineItem{Name: "widget", Quantity: 3}
	v := validation.Nested("item", item)
	if err := v.Validate(); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNestedValidateInnerFieldFails(t *testing.T) {
	item := LineItem{Name: "", Quantity: 3}
	v := validation.Nested("item", item)
	if err := v.Validate(); err == nil {
		t.Errorf("Error expected for empty name")
	}
}

func TestNestedValidateOuterRuleFails(t *testing.T) {
	// outer rule always fails; inner struct is valid
	item := LineItem{Name: "widget", Quantity: 3}
	reject := func(LineItem) *types.ValidationError {
		return types.NewValidationError("rejected")
	}
	v := validation.Nested("item", item, reject)
	if err := v.Validate(); err == nil {
		t.Errorf("Error expected from outer rule")
	}
}

func TestNestedValidateOuterRuleCheckedFirst(t *testing.T) {
	// outer rule fails; inner struct is also invalid — outer should surface first
	item := LineItem{Name: "", Quantity: -1}
	outerFired := false
	outer := func(LineItem) *types.ValidationError {
		outerFired = true
		return types.NewValidationError("outer failed")
	}
	v := validation.Nested("item", item, outer)
	err := v.Validate()
	if !outerFired {
		t.Errorf("Expected outer rule to be called")
	}
	if err == nil {
		t.Errorf("Error expected")
	}
}

func TestNestedValidateAllNoErrors(t *testing.T) {
	item := LineItem{Name: "widget", Quantity: 3}
	v := validation.Nested("item", item)
	errs := v.ValidateAll()
	if len(errs) != 0 {
		t.Errorf("No errors expected, got: %v", errs)
	}
}

func TestNestedValidateAllCollectsOuterAndInnerErrors(t *testing.T) {
	// outer rule fails + both inner fields fail
	item := LineItem{Name: "", Quantity: -1}
	reject := func(LineItem) *types.ValidationError {
		return types.NewValidationError("rejected")
	}
	v := validation.Nested("item", item, reject)
	errs := v.ValidateAll()
	if len(errs) < 3 {
		t.Errorf("Expected at least 3 errors (1 outer + 2 inner), got: %d", len(errs))
	}
}

func TestNestedValidateAllOnlyInnerErrors(t *testing.T) {
	// no outer rules; both inner fields fail
	item := LineItem{Name: "", Quantity: -1}
	v := validation.Nested("item", item)
	errs := v.ValidateAll()
	if len(errs) != 2 {
		t.Errorf("Expected 2 inner errors, got: %d", len(errs))
	}
}

type Invoice struct {
	ID   int
	Item LineItem
}

func (i Invoice) Validation() validation.Validation {
	return validation.NewValidations("invoice",
		validation.Value("ID", i.ID, rule.Positive[int]()),
		validation.Nested("item", i.Item),
	)
}

func TestNestedInsideNewValidations(t *testing.T) {
	inv := Invoice{ID: 1, Item: LineItem{Name: "bolt", Quantity: 10}}
	if err := inv.Validation().Validate(); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNestedInsideNewValidationsInnerFails(t *testing.T) {
	inv := Invoice{ID: 1, Item: LineItem{Name: "", Quantity: 10}}
	if err := inv.Validation().Validate(); err == nil {
		t.Errorf("Error expected for invalid nested item")
	}
}

func TestNestedValidateSkipsNilOuterRule(t *testing.T) {
	item := LineItem{Name: "widget", Quantity: 3}
	v := validation.Nested("item", item, nil)
	if err := v.Validate(); err != nil {
		t.Errorf("No errors expected when nil rule is skipped, got: %v", err)
	}
}

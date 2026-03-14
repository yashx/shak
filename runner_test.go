package shak

import (
	"testing"

	"github.com/yashx/shak/rule"
	"github.com/yashx/shak/types"
	"github.com/yashx/shak/validation"
)

// --- Validate ---

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

func TestValidateSkipsNilRule(t *testing.T) {
	err := Validate(10, nil, rule.Max[int](100))
	if err != nil {
		t.Errorf("No errors expected when nil rule is skipped, got: %v", err)
	}
}

func TestValidateStopsAtFirstError(t *testing.T) {
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
	Validate(0, r1, r2)
	if !firstRuleFired {
		t.Errorf("Expected first rule to be called")
	}
	if secondRuleFired {
		t.Errorf("Expected second rule to be skipped after first error")
	}
}

// --- Satisfies ---

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

// --- ValidateAll ---

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

func TestValidateAllSkipsNilRule(t *testing.T) {
	errs := ValidateAll(10, nil, rule.Max[int](100))
	if len(errs) != 0 {
		t.Errorf("No errors expected when nil rule is skipped, got: %v", errs)
	}
}

// --- RunValidation / RunFullValidation — named field paths ---

func TestRunValidationFieldPath(t *testing.T) {
	err := RunValidation(validation.Value("age", -1, rule.Min[int](0)))
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Path != "age" {
		t.Errorf("Expected path 'age', got: %v", err.Path)
	}
}

func TestRunFullValidationFieldPath(t *testing.T) {
	r1 := types.Rule[int](func(v int) *types.ValidationError {
		return types.NewValidationError("bad value")
	})
	errs := RunFullValidation(validation.Value("count", 0, r1))
	if len(errs) != 1 {
		t.Fatalf("Expected 1 error, got: %d", len(errs))
	}
	if errs[0].Path != "count" {
		t.Errorf("Expected path 'count', got: %v", errs[0].Path)
	}
}

func TestRunValidationSkipsNilValidation(t *testing.T) {
	type S struct{ X int }
	s := S{X: 5}
	v := validation.NewValidations(
		nil,
		validation.Value("X", s.X, rule.Max[int](10)),
	)
	if err := RunValidation(v); err != nil {
		t.Errorf("No errors expected when nil validation is skipped, got: %v", err)
	}
}

// --- RunValidation / RunFullValidation — struct validation via Validation() ---

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

// --- SatisfiesValidation ---

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

// --- Nested validation — direct Nested() calls ---

type lineItem struct {
	Name     string
	Quantity int
}

func (l lineItem) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("name", l.Name, rule.IsNotZeroValue[string]()),
		validation.Value("quantity", l.Quantity, rule.Positive[int]()),
	)
}

type invoice struct {
	ID   int
	Item lineItem
}

func (i invoice) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("ID", i.ID, rule.Positive[int]()),
		validation.Nested("item", i.Item),
	)
}

func TestRunValidationDirectNestedNoError(t *testing.T) {
	item := lineItem{Name: "widget", Quantity: 3}
	if err := RunValidation(validation.Nested("item", item)); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestRunValidationDirectNestedInnerFieldFails(t *testing.T) {
	item := lineItem{Name: "", Quantity: 3}
	if err := RunValidation(validation.Nested("item", item)); err == nil {
		t.Errorf("Error expected for empty name")
	}
}

func TestRunValidationDirectNestedOuterRuleCheckedFirst(t *testing.T) {
	item := lineItem{Name: "", Quantity: -1}
	outerFired := false
	outer := func(lineItem) *types.ValidationError {
		outerFired = true
		return types.NewValidationError("outer failed")
	}
	err := RunValidation(validation.Nested("item", item, outer))
	if !outerFired {
		t.Errorf("Expected outer rule to be called")
	}
	if err == nil {
		t.Errorf("Error expected")
	}
}

func TestRunFullValidationDirectNestedOnlyInnerErrors(t *testing.T) {
	item := lineItem{Name: "", Quantity: -1}
	errs := RunFullValidation(validation.Nested("item", item))
	if len(errs) != 2 {
		t.Errorf("Expected 2 inner errors, got: %d", len(errs))
	}
}

func TestRunFullValidationDirectNestedCollectsOuterAndInnerErrors(t *testing.T) {
	item := lineItem{Name: "", Quantity: -1}
	reject := func(lineItem) *types.ValidationError {
		return types.NewValidationError("rejected")
	}
	errs := RunFullValidation(validation.Nested("item", item, reject))
	if len(errs) < 3 {
		t.Errorf("Expected at least 3 errors (1 outer + 2 inner), got: %d", len(errs))
	}
}

func TestRunValidationDirectNestedSkipsNilOuterRule(t *testing.T) {
	item := lineItem{Name: "widget", Quantity: 3}
	if err := RunValidation(validation.Nested("item", item, nil)); err != nil {
		t.Errorf("No errors expected when nil rule is skipped, got: %v", err)
	}
}

// --- Nested validation — via struct Validation() ---

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

// --- Nested validation — path assertions ---

func TestRunValidationNestedInnerFieldPath(t *testing.T) {
	item := lineItem{Name: "", Quantity: 3}
	err := RunValidation(validation.Nested("item", item))
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Path != "item.name" {
		t.Errorf("Expected path 'item.name', got: %v", err.Path)
	}
}

func TestRunValidationNestedInnerSecondFieldPath(t *testing.T) {
	item := lineItem{Name: "bolt", Quantity: -1}
	err := RunValidation(validation.Nested("item", item))
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Path != "item.quantity" {
		t.Errorf("Expected path 'item.quantity', got: %v", err.Path)
	}
}

func TestRunValidationNestedOuterRulePath(t *testing.T) {
	item := lineItem{Name: "widget", Quantity: 3}
	reject := func(lineItem) *types.ValidationError {
		return types.NewValidationError("rejected")
	}
	err := RunValidation(validation.Nested("item", item, reject))
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Path != "item" {
		t.Errorf("Expected path 'item', got: %v", err.Path)
	}
}

func TestRunFullValidationNestedInnerPaths(t *testing.T) {
	item := lineItem{Name: "", Quantity: -1}
	errs := RunFullValidation(validation.Nested("item", item))
	if len(errs) != 2 {
		t.Fatalf("Expected 2 errors, got: %d", len(errs))
	}
	paths := map[string]bool{}
	for _, e := range errs {
		paths[e.Path] = true
	}
	if !paths["item.name"] {
		t.Errorf("Expected error at path 'item.name', got paths: %v", paths)
	}
	if !paths["item.quantity"] {
		t.Errorf("Expected error at path 'item.quantity', got paths: %v", paths)
	}
}

func TestRunFullValidationNestedOuterAndInnerPaths(t *testing.T) {
	item := lineItem{Name: "", Quantity: -1}
	reject := func(lineItem) *types.ValidationError {
		return types.NewValidationError("rejected")
	}
	errs := RunFullValidation(validation.Nested("item", item, reject))
	if len(errs) != 3 {
		t.Fatalf("Expected 3 errors, got: %d", len(errs))
	}
	paths := map[string]bool{}
	for _, e := range errs {
		paths[e.Path] = true
	}
	if !paths["item"] {
		t.Errorf("Expected error at path 'item' (outer rule), got paths: %v", paths)
	}
	if !paths["item.name"] {
		t.Errorf("Expected error at path 'item.name', got paths: %v", paths)
	}
	if !paths["item.quantity"] {
		t.Errorf("Expected error at path 'item.quantity', got paths: %v", paths)
	}
}

func TestRunValidationNestedInsideNewValidationsInnerPath(t *testing.T) {
	inv := invoice{ID: 1, Item: lineItem{Name: "", Quantity: 10}}
	err := RunValidation(inv.Validation())
	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
	if err.Path != "item.name" {
		t.Errorf("Expected path 'item.name', got: %v", err.Path)
	}
}

func TestRunFullValidationNestedInsideNewValidationsAllPaths(t *testing.T) {
	inv := invoice{ID: -1, Item: lineItem{Name: "", Quantity: -1}}
	errs := RunFullValidation(inv.Validation())
	if len(errs) != 3 {
		t.Fatalf("Expected 3 errors, got: %d", len(errs))
	}
	paths := map[string]bool{}
	for _, e := range errs {
		paths[e.Path] = true
	}
	if !paths["ID"] {
		t.Errorf("Expected error at path 'ID', got paths: %v", paths)
	}
	if !paths["item.name"] {
		t.Errorf("Expected error at path 'item.name', got paths: %v", paths)
	}
	if !paths["item.quantity"] {
		t.Errorf("Expected error at path 'item.quantity', got paths: %v", paths)
	}
}

// --- Validatable interface — compile-time check ---

func TestValidatableImplementsInterface(t *testing.T) {
	var _ validation.Validatable = address{}
	var _ validation.Validatable = person{}
	var _ validation.Validatable = lineItem{}
	var _ validation.Validatable = invoice{}
}

// --- Custom rules on nested ---
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

// --- Value vs Nested ---
//
// Using Value on a Validatable type does NOT run its Validation().
// Only the rules you explicitly pass are checked.

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
	s := shipment{TrackingID: "TRK-001", Drop: address{Street: "", City: ""}}
	if err := RunValidation(s.Validation()); err == nil {
		t.Errorf("Error expected: cityNotBlank rule should fail")
	}
}

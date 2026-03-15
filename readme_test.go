package shak_test

import (
	"testing"

	shak "github.com/yashx/shak"
	"github.com/yashx/shak/rule"
	"github.com/yashx/shak/types"
	"github.com/yashx/shak/validation"
)

// ---- Types from README examples ----

type Address struct {
	Street string
	City   string
}

func (a Address) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("street", a.Street, rule.NotBlank[string](), rule.MinLengthString[string](5)),
		validation.Value("city", a.City, rule.NotBlank[string]()),
	)
}

// BasicOrder is the Order from the basic Nested Structs example.
type BasicOrder struct {
	ID      int
	Address Address
}

func (o BasicOrder) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("ID", o.ID, rule.Positive[int]()),
		validation.Nested("address", o.Address),
	)
}

// ServicedOrder is the Order from the custom outer rule example.
type ServicedOrder struct {
	ID      int
	Address Address
}

var servicedCities = []string{"Austin", "Denver", "Portland"}

func mustBeServicedCity(a Address) *types.ValidationError {
	for _, c := range servicedCities {
		if a.City == c {
			return nil
		}
	}
	return types.NewValidationError("delivery not available in %q", a.City)
}

func (o ServicedOrder) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("ID", o.ID, rule.Positive[int]()),
		validation.Nested("address", o.Address, mustBeServicedCity),
	)
}

type Product struct {
	Name  string
	Price float64
}

func (p Product) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("name", p.Name, rule.NotBlank[string]()),
		validation.Value("price", p.Price, rule.Positive[float64]()),
	)
}

type Warehouse struct {
	Name     string
	Sections []map[string]int
}

func (w Warehouse) Validation() validation.Validation {
	var section map[string]int
	return validation.NewValidations(
		validation.Value("name", w.Name, rule.NotBlank[string]()),
		validation.Value("sections", w.Sections,
			rule.ForEach(&w.Sections,
				rule.ForEachValue(&section, rule.Min(0)),
			),
		),
	)
}

type Config struct {
	MaxRetries *int
}

func (c Config) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("maxRetries", c.MaxRetries, rule.Ref(rule.Min(0), rule.Max(10))),
	)
}

type GiftOrder struct {
	IsGift      bool
	GiftMessage string
}

func (o GiftOrder) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("giftMessage", o.GiftMessage,
			rule.If(o.IsGift, rule.NotBlank[string]()),
		),
	)
}

type TierOrder struct {
	IsFree bool
	Items  []string
}

func (o TierOrder) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("items", o.Items,
			rule.If(o.IsFree, rule.MaxLengthSlice(&o.Items, 5)).
				Else(rule.MaxLengthSlice(&o.Items, 100)),
		),
	)
}

type PersonWithAge struct {
	Age int
}

func (p PersonWithAge) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("age", p.Age,
			rule.Range(0, 120).WithMessage("age must be between 0 and 120"),
		),
	)
}

type NestedOrder struct {
	ID      int
	Address Address
}

func (o NestedOrder) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("ID", o.ID, rule.Positive[int]()),
		validation.Nested("address", o.Address),
	)
}

type ValueOrder struct {
	ID      int
	Address Address
}

func (o ValueOrder) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("ID", o.ID, rule.Positive[int]()),
		validation.Value("address", o.Address, func(a Address) *types.ValidationError {
			if a.City == "" {
				return types.NewValidationError("city is required for delivery")
			}
			return nil
		}),
	)
}

func Divisible(n int) types.Rule[int] {
	return func(v int) *types.ValidationError {
		if v%n == 0 {
			return nil
		}
		return types.NewValidationError("%d is not divisible by %d", v, n)
	}
}

// ---- Tests ----

func TestReadmeValidate(t *testing.T) {
	age := 150
	err := shak.Validate(age, rule.Min(0), rule.Max(120))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "value 150 is greater than maximum 120"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeValidateAll(t *testing.T) {
	score := 150
	errs := shak.ValidateAll(score, rule.Min(200), rule.Max(100))
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
	want := []string{
		"value 150 is less than minimum 200",
		"value 150 is greater than maximum 100",
	}
	for i, e := range errs {
		if e.Error() != want[i] {
			t.Errorf("[%d] got %q, want %q", i, e.Error(), want[i])
		}
	}
}

func TestReadmeRunValidation(t *testing.T) {
	a := Address{Street: "123", City: ""}
	err := shak.RunValidation(a)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "street: length is 3 bytes, expected at least 5"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeRunFullValidation(t *testing.T) {
	a := Address{Street: "123", City: ""}
	errs := shak.RunFullValidation(a)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
	want := []string{
		"street: length is 3 bytes, expected at least 5",
		"city: value must not be blank",
	}
	for i, e := range errs {
		if e.Error() != want[i] {
			t.Errorf("[%d] got %q, want %q", i, e.Error(), want[i])
		}
	}
}

func TestReadmeNestedStruct(t *testing.T) {
	o := BasicOrder{ID: 1, Address: Address{Street: "", City: "Austin"}}
	err := shak.RunValidation(o)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "address.street: value must not be blank"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeNestedOuterRule(t *testing.T) {
	o := ServicedOrder{ID: 1, Address: Address{Street: "1 Oak Ave", City: "Miami"}}
	err := shak.RunValidation(o)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := `address: delivery not available in "Miami"`
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeSatisfies(t *testing.T) {
	age := 10
	if !shak.Satisfies(age, rule.Min(0), rule.Max(120)) {
		t.Error("expected Satisfies to return true")
	}
}

func TestReadmeSatisfiesValidation(t *testing.T) {
	a := Address{Street: "1 Main St", City: "Austin"}
	if !shak.SatisfiesValidation(a) {
		t.Error("expected SatisfiesValidation to return true")
	}
}

func TestReadmeForEachValue(t *testing.T) {
	scores := map[string]int{"alice": 95, "bob": -1}
	err := shak.Validate(scores, rule.ForEachValue(&scores, rule.Min(0)))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "[bob]: value -1 is less than minimum 0"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeForEachKey(t *testing.T) {
	scores := map[string]int{"alice": 95, "bob": -1}
	err := shak.Validate(scores, rule.ForEachKey(&scores, rule.NotBlank[string]()))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestReadmeForKey(t *testing.T) {
	scores := map[string]int{"alice": 95, "bob": -1}
	err := shak.Validate(scores, rule.ForKey(&scores, "alice", rule.Max(100)))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestReadmeForEach(t *testing.T) {
	prices := []float64{9.99, -1.00, 49.99}
	err := shak.Validate(prices, rule.ForEach(&prices, rule.Positive[float64]()))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "[1]: value -1 is not positive"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeMinLengthSlice(t *testing.T) {
	prices := []float64{9.99, -1.00, 49.99}
	err := shak.Validate(prices, rule.MinLengthSlice(&prices, 1))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestReadmeUnique(t *testing.T) {
	tags := []string{"go", "go", "validation"}
	err := shak.Validate(tags, rule.Unique(&tags))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "duplicate value go found"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeIsValidSlice(t *testing.T) {
	products := []Product{
		{Name: "bolt", Price: 1.99},
		{Name: "", Price: -5},
	}
	err := shak.Validate(products, rule.ForEach(&products, rule.IsValid[Product]()))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "[1].name: value must not be blank"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeIsValidMap(t *testing.T) {
	catalog := map[string]Product{
		"A1": {Name: "bolt", Price: 1.99},
		"A2": {Name: "", Price: -5},
	}
	err := shak.Validate(catalog, rule.ForEachValue(&catalog, rule.IsValid[Product]()))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "[A2].name: value must not be blank"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeValidationErrorFields(t *testing.T) {
	a := Address{Street: "", City: "Austin"}
	err := shak.RunValidation(a)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Path != "street" {
		t.Errorf("Path: got %q, want %q", err.Path, "street")
	}
	if err.Message != "value must not be blank" {
		t.Errorf("Message: got %q, want %q", err.Message, "value must not be blank")
	}
	if err.Error() != "street: value must not be blank" {
		t.Errorf("Error(): got %q, want %q", err.Error(), "street: value must not be blank")
	}
}

func TestReadmeWarehouseErrorPath(t *testing.T) {
	w := Warehouse{
		Name: "Main",
		Sections: []map[string]int{
			{"bolts": 100, "nuts": 50},
			{"bolts": -3, "nuts": 200},
		},
	}
	err := shak.RunValidation(w)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "sections[1][bolts]: value -3 is less than minimum 0"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeDivisible(t *testing.T) {
	err := shak.Validate(10, Divisible(3))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "10 is not divisible by 3"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeAnd(t *testing.T) {
	validName := rule.And(rule.NotBlank[string](), rule.MinLengthString[string](3))
	err := shak.Validate("  ", validName)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "value must not be blank"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeOr(t *testing.T) {
	validInput := rule.Or(rule.LengthString[string](3), rule.MinLengthString[string](20))
	err := shak.Validate("hello", validInput)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "length is 5 bytes, expected at least 20"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeWithMessage(t *testing.T) {
	err := shak.Validate("", rule.NotBlank[string]().WithMessage("username is required"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "username is required"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeWithMessageOnRange(t *testing.T) {
	err := shak.RunValidation(PersonWithAge{Age: 150})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "age: age must be between 0 and 120"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeIfTrue(t *testing.T) {
	o := GiftOrder{IsGift: true, GiftMessage: ""}
	err := shak.RunValidation(o)
	if err == nil {
		t.Fatal("expected error: GiftMessage required when IsGift is true")
	}
	want := "giftMessage: value must not be blank"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeIfFalse(t *testing.T) {
	o := GiftOrder{IsGift: false, GiftMessage: ""}
	if err := shak.RunValidation(o); err != nil {
		t.Errorf("expected no error when IsGift is false, got %v", err)
	}
}

func TestReadmeIfElseFreeTier(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e", "f"} // 6 items
	o := TierOrder{IsFree: true, Items: items}
	if err := shak.RunValidation(o); err == nil {
		t.Fatal("expected error for free tier with 6 items")
	}
}

func TestReadmeIfElsePaidTier(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e", "f"} // 6 items
	o := TierOrder{IsFree: false, Items: items}
	if err := shak.RunValidation(o); err != nil {
		t.Errorf("expected no error for paid tier with 6 items, got %v", err)
	}
}

func TestReadmeRefNilPointer(t *testing.T) {
	c := Config{MaxRetries: nil}
	err := shak.RunValidation(c)
	if err == nil {
		t.Fatal("expected error for nil MaxRetries")
	}
	want := "maxRetries: value is nil"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeRefInvalidValue(t *testing.T) {
	n := -1
	c := Config{MaxRetries: &n}
	err := shak.RunValidation(c)
	if err == nil {
		t.Fatal("expected error for negative MaxRetries")
	}
	want := "maxRetries: value -1 is less than minimum 0"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestReadmeDeref(t *testing.T) {
	ptrRule := rule.Ref(rule.Positive[int]())
	err := shak.Validate(-1, rule.Deref(ptrRule))
	if err == nil {
		t.Fatal("expected error for non-positive value")
	}
	want := "value -1 is not positive"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

// TestReadmeValueVsNested demonstrates that validation.Nested calls the field's
// own Validation() while validation.Value does not.
func TestReadmeValueVsNested(t *testing.T) {
	// Nested: Address.Validation() runs — street and city are both checked.
	o := NestedOrder{ID: 1, Address: Address{Street: "", City: ""}}
	errs := shak.RunFullValidation(o)
	if len(errs) == 0 {
		t.Fatal("expected errors from Nested, got none")
	}
	// Address.Validation() enforces street NotBlank+MinLength and city NotBlank.
	found := false
	for _, e := range errs {
		if e.Error() == "address.street: value must not be blank" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected address.street error from Nested, got %v", errs)
	}

	// Value: only the explicit rule runs — Address.Validation() is ignored.
	o2 := ValueOrder{ID: 1, Address: Address{Street: "", City: ""}}
	err := shak.RunValidation(o2)
	if err == nil {
		t.Fatal("expected error from Value rule, got nil")
	}
	want := "address: city is required for delivery"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
	// Street is empty but Address.Validation() was NOT called, so no street error.
	errs2 := shak.RunFullValidation(o2)
	for _, e := range errs2 {
		if e.Path == "address.street" {
			t.Errorf("unexpected address.street error: Value should not call Address.Validation()")
		}
	}
}

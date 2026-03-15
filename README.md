# shak

**shak** (pronounced like "shuck") is a Hindi word (शक) that means *doubt*.

## Description

shak is a Go validation library built on **generics and higher-order functions** — no reflection, no struct tags, no magic. Rules are plain functions, and the type system enforces correctness at compile time.

It has the following features:

- **Type-safe by design** — powered by Go generics, not reflection.
- Use normal programming constructs rather than struct tags to specify validation rules.
- Validate data of different types — primitives, strings, slices, maps, structs, pointers, and more.
- Validate custom types by implementing the `Validatable` interface.
- Compose rules with `And`, `Or`, `If`, and `Else` for expressive conditional logic.
- Collect all errors at once or stop at the first failure.
- Structured errors with field paths for precise error reporting.
- Rich set of built-in rules, and easy to add custom ones.

## Installation

```
go get github.com/yashx/shak
```

## Requirements

Go 1.21 or above.

## A Note on LLM Usage

The design of this library — its architecture, core abstractions, and the initial set of validation rules — was written by me. LLMs were used to speed up the parts of development I don't find interesting and didn't want to spend my weekends on:

- **Tests** — unit and integration test coverage
- **Documentation** — doc comments and this README
- **Code review** — catching bugs and inconsistencies
- **Built-in rules** — filling out common use cases beyond the initial set

## Getting Started

shak is built around two concepts: **rules** and **validations**. A rule is a function that checks a single value. A validation groups rules and fields together and can be executed to collect errors.

- Use `shak.Validate()` / `shak.ValidateAll()` to validate a single value against rules directly.
- Use `shak.RunValidation()` / `shak.RunFullValidation()` to execute a structured validation built with `validation.Value()`, `validation.Nested()`, and `validation.NewValidations()`.
- Implement `Validatable` on your structs to make them self-describing.

### Validating a Simple Value

Use `shak.Validate()` to validate a single value against one or more rules:

```go
package main

import (
    "fmt"

    "github.com/yashx/shak"
    "github.com/yashx/shak/rule"
)

func main() {
    age := 150
    err := shak.Validate(age,
        rule.Min(0),   // must be >= 0
        rule.Max(120), // must be <= 120
    )
    fmt.Println(err)
    // Output:
    // value 150 is greater than maximum 120
}
```

Rules run in order. The first failure is returned and the rest are skipped. Use `shak.ValidateAll()` to collect every error instead:

```go
score := 150
errs := shak.ValidateAll(score,
    rule.Min(200), // must be >= 200
    rule.Max(100), // must be <= 100
)
for _, err := range errs {
    fmt.Println(err)
}
// Output:
// value 150 is less than minimum 200
// value 150 is greater than maximum 100
```

### Validating a Struct

To validate a struct, implement the `Validatable` interface by defining a `Validation()` method that describes the rules for each field. Then pass it to `shak.RunValidation()` or `shak.RunFullValidation()`:

```go
type Address struct {
    Street string
    City   string
}

func (a Address) Validation() validation.Validation {
    return validation.NewValidations(
        validation.Value("street", a.Street, rule.NotBlank[string](), rule.MinLengthString[string](5)),
        validation.Value("city",   a.City,   rule.NotBlank[string]()),
    )
}

a := Address{Street: "123", City: ""}

err := shak.RunValidation(a)
fmt.Println(err)
// Output:
// street: length is 3 bytes, expected at least 5

errs := shak.RunFullValidation(a)
for _, err := range errs {
    fmt.Println(err)
}
// Output:
// street: length is 3 bytes, expected at least 5
// city: value must not be blank
```

Fields are validated in the order they are declared. `RunValidation` stops at the first failure; `RunFullValidation` collects errors across all fields.

### Nested Structs

When a struct field is itself a `Validatable`, use `validation.Nested()` to include its own validation rules. Errors from the inner struct are automatically prefixed with the field name:

```go
type Address struct {
    Street string
    City   string
}

func (a Address) Validation() validation.Validation {
    return validation.NewValidations(
        validation.Value("street", a.Street, rule.NotBlank[string]()),
        validation.Value("city",   a.City,   rule.NotBlank[string]()),
    )
}

type Order struct {
    ID      int
    Address Address
}

func (o Order) Validation() validation.Validation {
    return validation.NewValidations(
        validation.Value("ID",      o.ID,      rule.Positive[int]()),
        validation.Nested("address", o.Address),
    )
}

o := Order{ID: 1, Address: Address{Street: "", City: "Austin"}}

err := shak.RunValidation(o)
fmt.Println(err)
// address.street: value must not be blank
```

You can also pass outer rules to `Nested`. They run before the field's own `Validation()` and receive the whole struct as the value:

```go
validCities := []string{"Austin", "Denver", "Portland"}

func mustBeServicedCity(a Address) *types.ValidationError {
    for _, c := range validCities {
        if a.City == c {
            return nil
        }
    }
    return types.NewValidationError("delivery not available in %q", a.City)
}

func (o Order) Validation() validation.Validation {
    return validation.NewValidations(
        validation.Value("ID",       o.ID,      rule.Positive[int]()),
        validation.Nested("address", o.Address, mustBeServicedCity),
    )
}

o := Order{ID: 1, Address: Address{Street: "1 Oak Ave", City: "Miami"}}

err := shak.RunValidation(o)
fmt.Println(err)
// address: delivery not available in "Miami"
```

### `Value` vs `Nested` for a Validatable Field

When a field's type implements `Validatable`, you have two choices:

- **`validation.Nested`** — runs any rules you pass *and* calls the field's own `Validation()`. Use this when you want the field to be fully validated as it defines itself.
- **`validation.Value`** — runs only the rules you explicitly pass. The field's own `Validation()` is never called. Use this when you want to apply a different or reduced set of rules in a specific context.

```go
// Nested: Address.Validation() runs — street and city are both checked
validation.Nested("address", o.Address)

// Value: only your rule runs — Address.Validation() is ignored
validation.Value("address", o.Address, func(a Address) *types.ValidationError {
    if a.City == "" {
        return types.NewValidationError("city is required for delivery")
    }
    return nil
})
```

### Checking Validity Without Errors

When you only need a boolean result and don't care about the error details, use `shak.Satisfies()` or `shak.SatisfiesValidation()`:

```go
// for a single value
if shak.Satisfies(age, rule.Min(0), rule.Max(120)) {
    fmt.Println("age is valid")
}

// for a struct validation
if shak.SatisfiesValidation(address) {
    fmt.Println("address is valid")
}
```

### Validating Maps

shak provides rules for validating map keys, values, and specific entries. Map rules take a pointer to the map as the first argument to help the compiler infer types:

```go
scores := map[string]int{"alice": 95, "bob": -1}

// validate every value
err := shak.Validate(scores, rule.ForEachValue(&scores, rule.Min(0)))
// [bob]: value -1 is less than minimum 0

// validate every key
err = shak.Validate(scores, rule.ForEachKey(&scores, rule.NotBlank[string]()))

// validate a specific key
err = shak.Validate(scores, rule.ForKey(&scores, "alice", rule.Max(100)))
```


Other available map rules:

- `rule.LengthMap(&m, n)` — map must have exactly n entries
- `rule.MinLengthMap(&m, n)` — map must have at least n entries
- `rule.MaxLengthMap(&m, n)` — map must have at most n entries
- `rule.EqualMap(m2)` — map must equal m2
- `rule.EqualMapFunc(&m, m2, eq)` — map must equal m2 using a custom equality function

### Validating Slices

Like map rules, slice rules take a pointer to the slice as the first argument for type inference:

```go
prices := []float64{9.99, -1.00, 49.99}

// validate every element
err := shak.Validate(prices, rule.ForEach(&prices, rule.Positive[float64]()))
// [1]: value -1 is not positive

// check length
err = shak.Validate(prices, rule.MinLengthSlice(&prices, 1))

// check for duplicates
tags := []string{"go", "go", "validation"}
err = shak.Validate(tags, rule.Unique(&tags))
// duplicate value go found
```

Other available slice rules:

- `rule.LengthSlice(&s, n)` — slice must have exactly n elements
- `rule.MaxLengthSlice(&s, n)` — slice must have at most n elements
- `rule.ContainsElem(&s, v)` — slice must contain v
- `rule.NotContainsElem(&s, v)` — slice must not contain v
- `rule.EqualSlice(s2)` — slice must equal s2
- `rule.EqualSliceFunc(&s, s2, eq)` — slice must equal s2 using a custom equality function
- `rule.IsSorted(&s)` — slice must be sorted in ascending order
- `rule.IsSortedFunc(&s, cmp)` — slice must be sorted by a custom comparison function
- `rule.MatchBytes(re)` — byte slice must match a regular expression
- `rule.NotMatchBytes(re)` — byte slice must not match a regular expression

### Validating a Slice or Map of Validatable Items

Use `rule.IsValid[T]()` to treat each element's own `Validation()` as a rule. This composes naturally with `ForEach` and `ForEachValue`:

```go
type Product struct {
    Name  string
    Price float64
}

func (p Product) Validation() validation.Validation {
    return validation.NewValidations(
        validation.Value("name",  p.Name,  rule.NotBlank[string]()),
        validation.Value("price", p.Price, rule.Positive[float64]()),
    )
}
```

**Slice:**

```go
products := []Product{
    {Name: "bolt", Price: 1.99},
    {Name: "",     Price: -5},  // invalid
}

err := shak.Validate(products, rule.ForEach(&products, rule.IsValid[Product]()))
fmt.Println(err)
// [1].name: value must not be blank
```

**Map:**

```go
catalog := map[string]Product{
    "A1": {Name: "bolt", Price: 1.99},
    "A2": {Name: "",     Price: -5},  // invalid
}

err := shak.Validate(catalog, rule.ForEachValue(&catalog, rule.IsValid[Product]()))
fmt.Println(err)
// [A2].name: value must not be blank
```

### Validation Errors

All validation functions return `*types.ValidationError`, which carries two fields:

```go
type ValidationError struct {
    Path    string `json:"path"`
    Message string `json:"message"`
}
```

- **`Path`** — the dot-separated path to the failing field (e.g. `"address.street"`, `"items[1].price"`). Empty for top-level single-value validation.
- **`Message`** — a human-readable description of the failure.

`ValidationError` implements the `error` interface. Calling `.Error()` returns `"path: message"`, or just `"message"` when there is no path:

```go
err := shak.RunValidation(address)
fmt.Println(err)         // street: value must not be blank
fmt.Println(err.Path)    // street
fmt.Println(err.Message) // value must not be blank
```

`RunFullValidation` and `ValidateAll` return `[]*types.ValidationError`, one entry per failing rule.

### Error Paths in Nested Structures

Paths are composed automatically as errors bubble up through nested structures. Consider a struct with a slice of maps:

```go
type Warehouse struct {
    Name     string
    Sections []map[string]int // section name -> stock count
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

w := Warehouse{
    Name: "Main",
    Sections: []map[string]int{
        {"bolts": 100, "nuts": 50},
        {"bolts": -3, "nuts": 200}, // invalid: negative stock
    },
}

err := shak.RunValidation(w)
fmt.Println(err)
// sections[1][bolts]: value -3 is less than minimum 0
```

The path `sections[1][bolts]` tells you exactly where the failure is — field `sections`, second element of the slice, key `bolts` in the map.

### Custom Rules

A rule is any function with the signature `func(T) *types.ValidationError`. Return `nil` to pass, or `types.NewValidationError(...)` to fail.

```go
func Divisible(n int) types.Rule[int] {
    return func(v int) *types.ValidationError {
        if v%n == 0 {
            return nil
        }
        return types.NewValidationError("%d is not divisible by %d", v, n)
    }
}

err := shak.Validate(10, Divisible(3))
fmt.Println(err)
// 10 is not divisible by 3
```

### Rule Composition

Use `rule.And` and `rule.Or` to combine multiple rules into one.

`rule.And` passes only if **all** rules pass — it stops at the first failure:

```go
// name must be non-blank and at least 3 bytes — two concerns, one rule
validName := rule.And(rule.NotBlank[string](), rule.MinLengthString[string](3))

err := shak.Validate("  ", validName)
fmt.Println(err)
// value must not be blank
```

`rule.Or` passes if **any** rule passes — it returns the last error only if all fail:

```go
// accept either a short code or a full description
validInput := rule.Or(rule.LengthString[string](3), rule.MinLengthString[string](20))

err := shak.Validate("hello", validInput)
fmt.Println(err)
// length is 5 bytes, expected at least 20
```

Composed rules are just rules — they can be passed anywhere a `types.Rule[T]` is accepted.

### Custom Error Messages

Any rule's default error message can be overridden with `WithMessage`. The error path is preserved:

```go
err := shak.Validate("", rule.NotBlank[string]().WithMessage("username is required"))
fmt.Println(err)
// username is required
```

It works on any rule, including composed ones:

```go
validation.Value("age", p.Age,
    rule.Range(0, 120).WithMessage("age must be between 0 and 120"),
)
```

### Conditional Rules

`rule.If` applies a rule only when a condition is true. When the condition is false it returns `nil`, which is safely skipped by the validation layer:

```go
type Order struct {
    IsGift      bool
    GiftMessage string
}

func (o Order) Validation() validation.Validation {
    return validation.NewValidations(
        // GiftMessage is only required when IsGift is true
        validation.Value("giftMessage", o.GiftMessage,
            rule.If(o.IsGift, rule.NotBlank[string]()),
        ),
    )
}
```

`Else` is a method on a rule — it provides a fallback for when `If` returns `nil`, letting you branch on a condition:

```go
// free tier: max 5 items; paid tier: max 100 items
validation.Value("items", o.Items,
    rule.If(o.IsFree, rule.MaxLengthSlice(&o.Items, 5)).
        Else(rule.MaxLengthSlice(&o.Items, 100)),
)
```

### Reusing Rules Across Pointer and Value Types

When a field is a pointer, you can't directly pass rules written for the value type. `rule.Ref` adapts value rules to work on a pointer — and returns an error automatically if the pointer is nil:

```go
type Config struct {
    MaxRetries *int
}

func (c Config) Validation() validation.Validation {
    return validation.NewValidations(
        // reuse Min/Max directly — no need to rewrite them for *int
        validation.Value("maxRetries", c.MaxRetries, rule.Ref(rule.Min(0), rule.Max(10))),
    )
}
```

`rule.Deref` goes the other way — it adapts pointer rules to work on a plain value by taking its address:

```go
// ptrRule is a Rule[*int] — reuse it on a plain int field
ptrRule := rule.Ref(rule.Positive[int]())
validation.Value("count", c.Count, rule.Deref(ptrRule))
```

### Built-in Rules

#### Ordered (`int`, `float64`, `string`, and any `cmp.Ordered` type)

- `rule.Min(n)` — value must be ≥ n
- `rule.Max(n)` — value must be ≤ n
- `rule.Range(min, max)` — value must be between min and max inclusive
- `rule.NotInRange(min, max)` — value must be outside [min, max]
- `rule.Positive()` — value must be > 0
- `rule.PositiveOrZero()` — value must be ≥ 0
- `rule.Negative()` — value must be < 0
- `rule.NegativeOrZero()` — value must be ≤ 0
- `rule.LessThan(n)` — value must be < n
- `rule.LessThanEqual(n)` — value must be ≤ n
- `rule.GreaterThan(n)` — value must be > n
- `rule.GreaterThanEqual(n)` — value must be ≥ n

#### Comparable (any `comparable` type)

- `rule.Equal(v)` — value must equal v
- `rule.NotEqual(v)` — value must not equal v
- `rule.EqualToAny(v...)` — value must match one of the provided values
- `rule.NotEqualToAny(v...)` — value must not match any of the provided values
- `rule.IsZeroValue()` — value must be the zero value for its type
- `rule.IsNotZeroValue()` — value must not be the zero value for its type

#### String (`string` and any `~string` type)

- `rule.NotBlank()` — string must not be empty or whitespace-only
- `rule.IsBlank()` — string must be empty or whitespace-only
- `rule.LengthString(n)` — string must be exactly n bytes
- `rule.MinLengthString(n)` — string must be at least n bytes
- `rule.MaxLengthString(n)` — string must be at most n bytes
- `rule.LengthRune(n)` — string must contain exactly n Unicode code points
- `rule.MinLengthRune(n)` — string must contain at least n Unicode code points
- `rule.MaxLengthRune(n)` — string must contain at most n Unicode code points
- `rule.ContainsString(s)` — string must contain substr s
- `rule.NotContainsString(s)` — string must not contain substr s
- `rule.HasPrefix(s)` — string must start with s
- `rule.NotHasPrefix(s)` — string must not start with s
- `rule.HasSuffix(s)` — string must end with s
- `rule.NotHasSuffix(s)` — string must not end with s
- `rule.ContainsAnyChar(s)` — string must contain at least one rune from s
- `rule.NotContainsAnyChar(s)` — string must contain none of the runes in s
- `rule.ContainsRune(r)` — string must contain rune r
- `rule.NotContainsRune(r)` — string must not contain rune r
- `rule.MatchString(re)` — string must match the regular expression
- `rule.NotMatchString(re)` — string must not match the regular expression
- `rule.EqualIgnoreCase(s)` — string must equal s, ignoring case
- `rule.NotEqualIgnoreCase(s)` — string must not equal s, ignoring case

#### Boolean (`bool` and any `~bool` type)

- `rule.IsTrue()` — value must be true
- `rule.IsFalse()` — value must be false

#### Pointer

- `rule.IsNil()` — pointer must be nil
- `rule.IsNotNil()` — pointer must not be nil

#### Time (`time.Time`)

- `rule.Before(t)` — time must be before t
- `rule.After(t)` — time must be after t
- `rule.Between(min, max)` — time must be between min and max inclusive
- `rule.InThePast()` — time must be before now
- `rule.InThePastOrPresent()` — time must be before or equal to now
- `rule.InTheFuture()` — time must be after now
- `rule.InTheFutureOrPresent()` — time must be after or equal to now
- `rule.OnDate(t)` — time must fall on the same calendar date as t

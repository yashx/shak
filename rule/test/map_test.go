package test

import (
	"testing"

	"github.com/yashx/shak/rule"
)

// ForEachValue

func TestForEachValueValid(t *testing.T) {
	m := map[string]int{"a": 2, "b": 5}
	r := rule.ForEachValue(&m, rule.Min[int](1))
	err := r(m)
	if err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestForEachValueError(t *testing.T) {
	m := map[string]int{"a": 1, "b": 20}
	r := rule.ForEachValue(&m, rule.Min[int](10))
	err := r(m)
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestForEachValueEmptyMap(t *testing.T) {
	m := map[string]int{}
	r := rule.ForEachValue(&m, rule.Min[int](10))
	err := r(m)
	if err != nil {
		t.Errorf("No errors expected for empty map, got: %v", err)
	}
}

func TestForEachValueErrorContainsKey(t *testing.T) {
	m := map[string]int{"badkey": 1}
	r := rule.ForEachValue(&m, rule.Min[int](10))
	err := r(m)
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
	if err.Path != "[badkey]" {
		t.Errorf("Expected path [badkey], got: %v", err.Path)
	}
}

func TestForEachValueCustomType(t *testing.T) {
	type MyMap map[string]int
	m := MyMap{"x": 5}
	r := rule.ForEachValue(&m, rule.Min[int](1))
	err := r(m)
	if err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// ForKey

func TestForKeyValid(t *testing.T) {
	m := map[string]int{"age": 21}
	r := rule.ForKey(&m,"age", rule.Min[int](18))
	err := r(m)
	if err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestForKeyError(t *testing.T) {
	m := map[string]int{"age": 10}
	r := rule.ForKey(&m,"age", rule.Min[int](18))
	err := r(m)
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestForKeyErrorContainsKey(t *testing.T) {
	m := map[string]int{"age": 10}
	r := rule.ForKey(&m,"age", rule.Min[int](18))
	err := r(m)
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
	if err.Path != "[age]" {
		t.Errorf("Expected path [age], got: %v", err.Path)
	}
}

func TestForKeyMissingKey(t *testing.T) {
	m := map[string]int{}
	r := rule.ForKey(&m,"missing", rule.Min[int](1))
	err := r(m)
	if err == nil {
		t.Errorf("Error expected for missing key, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestForKeyMissingKeyErrorPath(t *testing.T) {
	m := map[string]int{}
	r := rule.ForKey(&m,"score", rule.Min[int](1))
	err := r(m)
	if err == nil {
		t.Fatalf("Error expected for missing key, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
	if err.Path != "[score]" {
		t.Errorf("Expected path [score], got: %v", err.Path)
	}
}

func TestForKeyZeroValueExplicitlySet(t *testing.T) {
	// key exists with zero value — should pass through to the rule, not be treated as missing
	m := map[string]int{"count": 0}
	r := rule.ForKey(&m,"count", rule.Min[int](0))
	if err := r(m); err != nil {
		t.Errorf("No errors expected for explicitly set zero value, got: %v", err)
	}
}

func TestForKeyZeroValueExplicitlySetFails(t *testing.T) {
	// key exists with zero value but fails the rule — should return rule error, not missing-key error
	m := map[string]int{"count": 0}
	r := rule.ForKey(&m,"count", rule.Positive[int]())
	err := r(m)
	if err == nil {
		t.Errorf("Error expected when explicitly set zero value fails rule")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestForKeyIntKey(t *testing.T) {
	m := map[int]string{1: "hello"}
	r := rule.ForKey(&m,1, rule.LengthString[string](5))
	err := r(m)
	if err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestForKeyCustomType(t *testing.T) {
	type MyMap map[string]int
	m := MyMap{"score": 100}
	r := rule.ForKey(&m,"score", rule.Min[int](0))
	err := r(m)
	if err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// LengthMap

func TestLengthMapValid(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	r := rule.LengthMap(&m, 2)
	err := r(m)
	if err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestLengthMapError(t *testing.T) {
	m := map[string]int{"a": 1}
	r := rule.LengthMap(&m, 3)
	err := r(m)
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestLengthMapEmpty(t *testing.T) {
	m := map[string]int{}
	r := rule.LengthMap(&m, 0)
	err := r(m)
	if err != nil {
		t.Errorf("No errors expected for empty map with length 0, got: %v", err)
	}
}

func TestLengthMapCustomType(t *testing.T) {
	type MyMap map[string]bool
	m := MyMap{"only": true}
	r := rule.LengthMap(&m, 1)
	err := r(m)
	if err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestForEachValueIsValid(t *testing.T) {
	m := map[string]product{
		"a": {Name: "Widget", Price: 9.99},
		"b": {Name: "Gadget", Price: 4.99},
	}
	r := rule.ForEachValue(&m, rule.IsValid[product]())
	if err := r(m); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestForEachValueIsValidOneInvalid(t *testing.T) {
	m := map[string]product{
		"a": {Name: "Widget", Price: 9.99},
		"b": {Name: "", Price: 4.99},
	}
	r := rule.ForEachValue(&m, rule.IsValid[product]())
	err := r(m)
	if err == nil {
		t.Errorf("Error expected for invalid product, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

// EqualMap

func TestEqualMapValid(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	r := rule.EqualMap(map[string]int{"a": 1, "b": 2})
	if err := r(m); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestEqualMapError(t *testing.T) {
	m := map[string]int{"a": 1}
	r := rule.EqualMap(map[string]int{"a": 2})
	err := r(m)
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestEqualMapEmpty(t *testing.T) {
	m := map[string]int{}
	r := rule.EqualMap(map[string]int{})
	if err := r(m); err != nil {
		t.Errorf("No errors expected for equal empty maps, got: %v", err)
	}
}

func TestEqualMapDifferentKeys(t *testing.T) {
	m := map[string]int{"a": 1}
	r := rule.EqualMap(map[string]int{"b": 1})
	err := r(m)
	if err == nil {
		t.Errorf("Error expected for maps with different keys, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestEqualMapCustomType(t *testing.T) {
	type MyMap map[string]int
	m := MyMap{"x": 10}
	r := rule.EqualMap(MyMap{"x": 10})
	if err := r(m); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// EqualMapFunc

func TestEqualMapFuncValid(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	r := rule.EqualMapFunc(&m, map[string]int{"a": 1, "b": 2}, func(a, b int) bool { return a == b })
	if err := r(m); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestEqualMapFuncError(t *testing.T) {
	m := map[string]int{"a": 1}
	r := rule.EqualMapFunc(&m, map[string]int{"a": 2}, func(a, b int) bool { return a == b })
	err := r(m)
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestEqualMapFuncCustomEquality(t *testing.T) {
	// Compare ints to floats, different value types
	m := map[string]int{"a": 1, "b": 2}
	r := rule.EqualMapFunc(&m, map[string]float64{"a": 1.0, "b": 2.0}, func(a int, b float64) bool { return float64(a) == b })
	if err := r(m); err != nil {
		t.Errorf("No errors expected with custom equality, got: %v", err)
	}
}

func TestEqualMapFuncCustomEqualityError(t *testing.T) {
	m := map[string]int{"a": 1}
	r := rule.EqualMapFunc(&m, map[string]float64{"a": 1.9}, func(a int, b float64) bool { return float64(a) == b })
	err := r(m)
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

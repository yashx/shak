package test

import (
	"regexp"
	"testing"

	"github.com/yashx/shak/rule"
)

// ForEach

func TestForEachIntValid(t *testing.T) {
	s := []int{2, 5, 10}
	r := rule.ForEach(&s,rule.Min[int](1))
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestForEachIntError(t *testing.T) {
	s := []int{5, -1, 10}
	r := rule.ForEach(&s,rule.Positive[int]())
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for slice containing negative value")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestForEachErrorContainsIndex(t *testing.T) {
	s := []int{5, -1, 10}
	r := rule.ForEach(&s,rule.Positive[int]())
	err := r(s)
	if err == nil {
		t.Fatalf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
	if err.Path != "[1]" {
		t.Errorf("Expected path [1], got: %v", err.Path)
	}
}

func TestForEachEmptySlice(t *testing.T) {
	s := []int{}
	r := rule.ForEach(&s,rule.Positive[int]())
	if err := r(s); err != nil {
		t.Errorf("No errors expected for empty slice, got: %v", err)
	}
}

func TestForEachStringValid(t *testing.T) {
	s := []string{"hello", "world"}
	r := rule.ForEach(&s,rule.IsNotZeroValue[string]())
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestForEachStringError(t *testing.T) {
	s := []string{"hello", "", "world"}
	r := rule.ForEach(&s,rule.IsNotZeroValue[string]())
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for slice containing empty string")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestForEachCustomType(t *testing.T) {
	type Score int
	s := []Score{10, 20, 30}
	r := rule.ForEach(&s,rule.Positive[Score]())
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestForEachStopsAtFirstError(t *testing.T) {
	s := []int{-1, -2, -3}
	r := rule.ForEach(&s,rule.Positive[int]())
	err := r(s)
	if err == nil {
		t.Fatalf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
	// should report the first failing index
	if err.Path != "[0]" {
		t.Errorf("Expected path [0] (first error), got: %v", err.Path)
	}
}

// LengthSlice

func TestLengthSliceIntValid(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.LengthSlice(&s, 3)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestLengthSliceIntError(t *testing.T) {
	s := []int{1, 2}
	r := rule.LengthSlice(&s, 5)
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for wrong length")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestLengthSliceEmpty(t *testing.T) {
	s := []string{}
	r := rule.LengthSlice(&s, 0)
	if err := r(s); err != nil {
		t.Errorf("No errors expected for empty slice with length 0, got: %v", err)
	}
}

func TestLengthSliceString(t *testing.T) {
	s := []string{"a", "b"}
	r := rule.LengthSlice(&s, 2)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestLengthSliceCustomType(t *testing.T) {
	type Tags []string
	s := Tags{"go", "generics"}
	r := rule.LengthSlice(&s, 2)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestLengthSliceCustomTypeError(t *testing.T) {
	type Tags []string
	s := Tags{"go"}
	r := rule.LengthSlice(&s, 3)
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for wrong length")
	} else {
		t.Logf("got error: %v", err)
	}
}

// MinLengthSlice

func TestMinLengthSliceValid(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.MinLengthSlice(&s, 2)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMinLengthSliceAtBound(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.MinLengthSlice(&s, 3)
	if err := r(s); err != nil {
		t.Errorf("No errors expected at exact min length, got: %v", err)
	}
}

func TestMinLengthSliceError(t *testing.T) {
	s := []int{1}
	r := rule.MinLengthSlice(&s, 3)
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for slice below min length")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestMinLengthSliceCustomType(t *testing.T) {
	type Tags []string
	s := Tags{"a", "b"}
	r := rule.MinLengthSlice(&s, 1)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// MaxLengthSlice

func TestMaxLengthSliceValid(t *testing.T) {
	s := []int{1, 2}
	r := rule.MaxLengthSlice(&s, 5)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMaxLengthSliceAtBound(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.MaxLengthSlice(&s, 3)
	if err := r(s); err != nil {
		t.Errorf("No errors expected at exact max length, got: %v", err)
	}
}

func TestMaxLengthSliceError(t *testing.T) {
	s := []int{1, 2, 3, 4}
	r := rule.MaxLengthSlice(&s, 3)
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for slice above max length")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestMaxLengthSliceCustomType(t *testing.T) {
	type Tags []string
	s := Tags{"a"}
	r := rule.MaxLengthSlice(&s, 5)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotContainsElem

func TestNotContainsElemValid(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.NotContainsElem(&s, 99)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotContainsElemError(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.NotContainsElem(&s, 2)
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for slice containing forbidden element")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotContainsElemEmpty(t *testing.T) {
	s := []int{}
	r := rule.NotContainsElem(&s, 1)
	if err := r(s); err != nil {
		t.Errorf("No errors expected for empty slice, got: %v", err)
	}
}

func TestNotContainsElemCustomType(t *testing.T) {
	type Role string
	s := []Role{"user", "editor"}
	r := rule.NotContainsElem(&s, Role("admin"))
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotMatchBytes

func TestNotMatchBytesValid(t *testing.T) {
	re := regexp.MustCompile(`\d`)
	r := rule.NotMatchBytes[[]byte](re)
	if err := r([]byte("hello")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotMatchBytesError(t *testing.T) {
	re := regexp.MustCompile(`\d`)
	r := rule.NotMatchBytes[[]byte](re)
	err := r([]byte("hello123"))
	if err == nil {
		t.Errorf("Error expected for bytes matching pattern")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotMatchBytesCustomType(t *testing.T) {
	type Payload []byte
	re := regexp.MustCompile(`^\x00`)
	r := rule.NotMatchBytes[Payload](re)
	if err := r(Payload("normal")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// Unique

func TestUniqueValid(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.Unique(&s)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestUniqueError(t *testing.T) {
	s := []int{1, 2, 2, 3}
	r := rule.Unique(&s)
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for slice with duplicate elements")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestUniqueEmpty(t *testing.T) {
	s := []int{}
	r := rule.Unique(&s)
	if err := r(s); err != nil {
		t.Errorf("No errors expected for empty slice, got: %v", err)
	}
}

func TestUniqueSingleElement(t *testing.T) {
	s := []int{42}
	r := rule.Unique(&s)
	if err := r(s); err != nil {
		t.Errorf("No errors expected for single-element slice, got: %v", err)
	}
}

func TestUniqueStrings(t *testing.T) {
	s := []string{"a", "b", "c"}
	r := rule.Unique(&s)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestUniqueStringsError(t *testing.T) {
	s := []string{"a", "b", "a"}
	r := rule.Unique(&s)
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for slice with duplicate strings")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestUniqueCustomType(t *testing.T) {
	type ID int
	s := []ID{1, 2, 3}
	r := rule.Unique(&s)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// Contains

func TestContainsValid(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.ContainsElem(&s, 2)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestContainsError(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.ContainsElem(&s, 5)
	err := r(s)
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestContainsEmptySlice(t *testing.T) {
	s := []int{}
	r := rule.ContainsElem(&s, 1)
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for empty slice, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestContainsCustomType(t *testing.T) {
	type Tag string
	s := []Tag{"go", "generics"}
	r := rule.ContainsElem(&s, Tag("go"))
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// EqualSlice

func TestEqualSliceValid(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.EqualSlice([]int{1, 2, 3})
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestEqualSliceError(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.EqualSlice([]int{1, 2, 4})
	err := r(s)
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestEqualSliceEmpty(t *testing.T) {
	s := []int{}
	r := rule.EqualSlice([]int{})
	if err := r(s); err != nil {
		t.Errorf("No errors expected for equal empty slices, got: %v", err)
	}
}

func TestEqualSliceDifferentLength(t *testing.T) {
	s := []int{1, 2}
	r := rule.EqualSlice([]int{1, 2, 3})
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for different lengths, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestEqualSliceCustomType(t *testing.T) {
	type Ids []int
	s := Ids{10, 20}
	r := rule.EqualSlice(Ids{10, 20})
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// EqualSliceFunc

func TestEqualSliceFuncValid(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.EqualSliceFunc(&s, []int{1, 2, 3}, func(a, b int) bool { return a == b })
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestEqualSliceFuncError(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.EqualSliceFunc(&s, []int{1, 2, 4}, func(a, b int) bool { return a == b })
	err := r(s)
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestEqualSliceFuncDifferentTypes(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.EqualSliceFunc(&s, []float64{1.0, 2.0, 3.0}, func(a int, b float64) bool { return float64(a) == b })
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestEqualSliceFuncDifferentTypesError(t *testing.T) {
	s := []int{1, 2, 3}
	r := rule.EqualSliceFunc(&s, []float64{1.0, 2.0, 3.5}, func(a int, b float64) bool { return float64(a) == b })
	err := r(s)
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

// IsSorted

func TestIsSortedValid(t *testing.T) {
	s := []int{1, 2, 3, 4}
	r := rule.IsSorted(&s)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIsSortedError(t *testing.T) {
	s := []int{3, 1, 2}
	r := rule.IsSorted(&s)
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for unsorted slice, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestIsSortedEmpty(t *testing.T) {
	s := []int{}
	r := rule.IsSorted(&s)
	if err := r(s); err != nil {
		t.Errorf("No errors expected for empty slice, got: %v", err)
	}
}

func TestIsSortedStrings(t *testing.T) {
	s := []string{"a", "b", "c"}
	r := rule.IsSorted(&s)
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// IsSortedFunc

func TestIsSortedFuncValid(t *testing.T) {
	s := []int{4, 3, 2, 1}
	r := rule.IsSortedFunc(&s, func(a, b int) int { return b - a })
	if err := r(s); err != nil {
		t.Errorf("No errors expected for descending sorted slice, got: %v", err)
	}
}

func TestIsSortedFuncError(t *testing.T) {
	s := []int{1, 3, 2}
	r := rule.IsSortedFunc(&s, func(a, b int) int { return b - a })
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for unsorted slice, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestIsSortedFuncCustomType(t *testing.T) {
	type Person struct{ Age int }
	s := []Person{{1}, {2}, {3}}
	r := rule.IsSortedFunc(&s, func(a, b Person) int { return a.Age - b.Age })
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// Match (works on ~[]byte)

func TestMatchBytesValid(t *testing.T) {
	re := regexp.MustCompile(`^\d+$`)
	r := rule.MatchBytes[[]byte](re)
	if err := r([]byte("12345")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMatchBytesError(t *testing.T) {
	re := regexp.MustCompile(`^\d+$`)
	r := rule.MatchBytes[[]byte](re)
	err := r([]byte("abc"))
	if err == nil {
		t.Errorf("Error expected for non-digit bytes")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestMatchCustomBytesType(t *testing.T) {
	type Payload []byte
	re := regexp.MustCompile(`^hello`)
	r := rule.MatchBytes[Payload](re)
	if err := r(Payload("hello world")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMatchCustomBytesTypeError(t *testing.T) {
	type Payload []byte
	re := regexp.MustCompile(`^hello`)
	r := rule.MatchBytes[Payload](re)
	err := r(Payload("goodbye"))
	if err == nil {
		t.Errorf("Error expected for non-matching payload")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestForEachIsValid(t *testing.T) {
	s := []product{
		{Name: "Widget", Price: 9.99},
		{Name: "Gadget", Price: 4.99},
	}
	r := rule.ForEach(&s, rule.IsValid[product]())
	if err := r(s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestForEachIsValidOneInvalid(t *testing.T) {
	s := []product{
		{Name: "Widget", Price: 9.99},
		{Name: "Gadget", Price: -1.0},
	}
	r := rule.ForEach(&s, rule.IsValid[product]())
	err := r(s)
	if err == nil {
		t.Errorf("Error expected for invalid product, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

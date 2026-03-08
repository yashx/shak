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
	if err := r(s); err == nil {
		t.Errorf("Error expected for slice containing negative value")
	}
}

func TestForEachErrorContainsIndex(t *testing.T) {
	s := []int{5, -1, 10}
	r := rule.ForEach(&s,rule.Positive[int]())
	err := r(s)
	if err == nil {
		t.Fatalf("Error expected, got nil")
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
	if err := r(s); err == nil {
		t.Errorf("Error expected for slice containing empty string")
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
	if err := r(s); err == nil {
		t.Errorf("Error expected for wrong length")
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
	if err := r(s); err == nil {
		t.Errorf("Error expected for wrong length")
	}
}

// Match (works on ~[]byte)

func TestMatchBytesValid(t *testing.T) {
	re := regexp.MustCompile(`^\d+$`)
	r := rule.Match[[]byte](re)
	if err := r([]byte("12345")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMatchBytesError(t *testing.T) {
	re := regexp.MustCompile(`^\d+$`)
	r := rule.Match[[]byte](re)
	if err := r([]byte("abc")); err == nil {
		t.Errorf("Error expected for non-digit bytes")
	}
}

func TestMatchCustomBytesType(t *testing.T) {
	type Payload []byte
	re := regexp.MustCompile(`^hello`)
	r := rule.Match[Payload](re)
	if err := r(Payload("hello world")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMatchCustomBytesTypeError(t *testing.T) {
	type Payload []byte
	re := regexp.MustCompile(`^hello`)
	r := rule.Match[Payload](re)
	if err := r(Payload("goodbye")); err == nil {
		t.Errorf("Error expected for non-matching payload")
	}
}

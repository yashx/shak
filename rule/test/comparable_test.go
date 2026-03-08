package test

import (
	"testing"

	"github.com/yashx/shak/rule"
)

// IsZeroValue

func TestIsZeroValueInt(t *testing.T) {
	r := rule.IsZeroValue[int]()
	if err := r(0); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIsZeroValueIntError(t *testing.T) {
	r := rule.IsZeroValue[int]()
	if err := r(1); err == nil {
		t.Errorf("Error expected for non-zero int")
	}
}

func TestIsZeroValueString(t *testing.T) {
	r := rule.IsZeroValue[string]()
	if err := r(""); err != nil {
		t.Errorf("No errors expected for empty string, got: %v", err)
	}
}

func TestIsZeroValueStringError(t *testing.T) {
	r := rule.IsZeroValue[string]()
	if err := r("hello"); err == nil {
		t.Errorf("Error expected for non-empty string")
	}
}

func TestIsZeroValueBool(t *testing.T) {
	r := rule.IsZeroValue[bool]()
	if err := r(false); err != nil {
		t.Errorf("No errors expected for false, got: %v", err)
	}
}

func TestIsZeroValueBoolError(t *testing.T) {
	r := rule.IsZeroValue[bool]()
	if err := r(true); err == nil {
		t.Errorf("Error expected for true")
	}
}

func TestIsZeroValueCustomType(t *testing.T) {
	type Status int
	r := rule.IsZeroValue[Status]()
	if err := r(Status(0)); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// IsNotZeroValue

func TestIsNotZeroValueInt(t *testing.T) {
	r := rule.IsNotZeroValue[int]()
	if err := r(42); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIsNotZeroValueIntError(t *testing.T) {
	r := rule.IsNotZeroValue[int]()
	if err := r(0); err == nil {
		t.Errorf("Error expected for zero int")
	}
}

func TestIsNotZeroValueString(t *testing.T) {
	r := rule.IsNotZeroValue[string]()
	if err := r("world"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIsNotZeroValueStringError(t *testing.T) {
	r := rule.IsNotZeroValue[string]()
	if err := r(""); err == nil {
		t.Errorf("Error expected for empty string")
	}
}

func TestIsNotZeroValueFloat64(t *testing.T) {
	r := rule.IsNotZeroValue[float64]()
	if err := r(3.14); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIsNotZeroValueFloat64Error(t *testing.T) {
	r := rule.IsNotZeroValue[float64]()
	if err := r(0.0); err == nil {
		t.Errorf("Error expected for zero float64")
	}
}

// Equal

func TestEqualInt(t *testing.T) {
	r := rule.Equal[int](10)
	if err := r(10); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestEqualIntError(t *testing.T) {
	r := rule.Equal[int](10)
	if err := r(11); err == nil {
		t.Errorf("Error expected for unequal int")
	}
}

func TestEqualString(t *testing.T) {
	r := rule.Equal[string]("go")
	if err := r("go"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestEqualStringError(t *testing.T) {
	r := rule.Equal[string]("go")
	if err := r("rust"); err == nil {
		t.Errorf("Error expected for unequal string")
	}
}

func TestEqualFloat64(t *testing.T) {
	r := rule.Equal[float64](1.5)
	if err := r(1.5); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotEqual

func TestNotEqualInt(t *testing.T) {
	r := rule.NotEqual[int](0)
	if err := r(1); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotEqualIntError(t *testing.T) {
	r := rule.NotEqual[int](0)
	if err := r(0); err == nil {
		t.Errorf("Error expected when value equals forbidden value")
	}
}

func TestNotEqualString(t *testing.T) {
	r := rule.NotEqual[string]("admin")
	if err := r("user"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotEqualStringError(t *testing.T) {
	r := rule.NotEqual[string]("admin")
	if err := r("admin"); err == nil {
		t.Errorf("Error expected when value equals forbidden string")
	}
}

func TestNotEqualBool(t *testing.T) {
	r := rule.NotEqual[bool](false)
	if err := r(true); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// EqualsAny

func TestEqualsAnyInt(t *testing.T) {
	r := rule.EqualsAny[int](1, 2, 3)
	if err := r(2); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestEqualsAnyIntError(t *testing.T) {
	r := rule.EqualsAny[int](1, 2, 3)
	if err := r(99); err == nil {
		t.Errorf("Error expected for value not in list")
	}
}

func TestEqualsAnyString(t *testing.T) {
	r := rule.EqualsAny[string]("red", "green", "blue")
	if err := r("green"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestEqualsAnyStringError(t *testing.T) {
	r := rule.EqualsAny[string]("red", "green", "blue")
	if err := r("purple"); err == nil {
		t.Errorf("Error expected for value not in list")
	}
}

func TestEqualsAnyCustomType(t *testing.T) {
	type Color string
	r := rule.EqualsAny[Color]("red", "green")
	if err := r(Color("red")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotEqualsAny

func TestNotEqualsAnyInt(t *testing.T) {
	r := rule.NotEqualsAny[int](0, -1)
	if err := r(5); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotEqualsAnyIntError(t *testing.T) {
	r := rule.NotEqualsAny[int](0, -1)
	if err := r(0); err == nil {
		t.Errorf("Error expected when value is in forbidden list")
	}
}

func TestNotEqualsAnyString(t *testing.T) {
	r := rule.NotEqualsAny[string]("root", "admin")
	if err := r("alice"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotEqualsAnyStringError(t *testing.T) {
	r := rule.NotEqualsAny[string]("root", "admin")
	if err := r("admin"); err == nil {
		t.Errorf("Error expected when value is in forbidden list")
	}
}

func TestNotEqualsAnyBool(t *testing.T) {
	r := rule.NotEqualsAny[bool](false)
	if err := r(true); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

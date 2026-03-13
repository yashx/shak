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
	err := r(1)
	if err == nil {
		t.Errorf("Error expected for non-zero int")
	} else {
		t.Logf("got error: %v", err)
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
	err := r("hello")
	if err == nil {
		t.Errorf("Error expected for non-empty string")
	} else {
		t.Logf("got error: %v", err)
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
	err := r(true)
	if err == nil {
		t.Errorf("Error expected for true")
	} else {
		t.Logf("got error: %v", err)
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
	err := r(0)
	if err == nil {
		t.Errorf("Error expected for zero int")
	} else {
		t.Logf("got error: %v", err)
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
	err := r("")
	if err == nil {
		t.Errorf("Error expected for empty string")
	} else {
		t.Logf("got error: %v", err)
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
	err := r(0.0)
	if err == nil {
		t.Errorf("Error expected for zero float64")
	} else {
		t.Logf("got error: %v", err)
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
	err := r(11)
	if err == nil {
		t.Errorf("Error expected for unequal int")
	} else {
		t.Logf("got error: %v", err)
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
	err := r("rust")
	if err == nil {
		t.Errorf("Error expected for unequal string")
	} else {
		t.Logf("got error: %v", err)
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
	err := r(0)
	if err == nil {
		t.Errorf("Error expected when value equals forbidden value")
	} else {
		t.Logf("got error: %v", err)
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
	err := r("admin")
	if err == nil {
		t.Errorf("Error expected when value equals forbidden string")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotEqualBool(t *testing.T) {
	r := rule.NotEqual[bool](false)
	if err := r(true); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// EqualToAny

func TestEqualToAnyInt(t *testing.T) {
	r := rule.EqualToAny[int](1, 2, 3)
	if err := r(2); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestEqualToAnyIntError(t *testing.T) {
	r := rule.EqualToAny[int](1, 2, 3)
	err := r(99)
	if err == nil {
		t.Errorf("Error expected for value not in list")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestEqualToAnyString(t *testing.T) {
	r := rule.EqualToAny[string]("red", "green", "blue")
	if err := r("green"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestEqualToAnyStringError(t *testing.T) {
	r := rule.EqualToAny[string]("red", "green", "blue")
	err := r("purple")
	if err == nil {
		t.Errorf("Error expected for value not in list")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestEqualToAnyCustomType(t *testing.T) {
	type Color string
	r := rule.EqualToAny[Color]("red", "green")
	if err := r(Color("red")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotEqualToAny

func TestNotEqualToAnyInt(t *testing.T) {
	r := rule.NotEqualToAny[int](0, -1)
	if err := r(5); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotEqualToAnyIntError(t *testing.T) {
	r := rule.NotEqualToAny[int](0, -1)
	err := r(0)
	if err == nil {
		t.Errorf("Error expected when value is in forbidden list")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotEqualToAnyString(t *testing.T) {
	r := rule.NotEqualToAny[string]("root", "admin")
	if err := r("alice"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotEqualToAnyStringError(t *testing.T) {
	r := rule.NotEqualToAny[string]("root", "admin")
	err := r("admin")
	if err == nil {
		t.Errorf("Error expected when value is in forbidden list")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotEqualToAnyBool(t *testing.T) {
	r := rule.NotEqualToAny[bool](false)
	if err := r(true); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

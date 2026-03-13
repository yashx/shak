package test

import (
	"testing"

	"github.com/yashx/shak/rule"
)

// IsNil

func TestIsNilIntNil(t *testing.T) {
	r := rule.IsNil[int]()
	var p *int = nil
	if err := r(p); err != nil {
		t.Errorf("No errors expected for nil pointer, got: %v", err)
	}
}

func TestIsNilIntError(t *testing.T) {
	r := rule.IsNil[int]()
	v := 42
	err := r(&v)
	if err == nil {
		t.Errorf("Error expected for non-nil pointer")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestIsNilStringNil(t *testing.T) {
	r := rule.IsNil[string]()
	var p *string = nil
	if err := r(p); err != nil {
		t.Errorf("No errors expected for nil string pointer, got: %v", err)
	}
}

func TestIsNilStringError(t *testing.T) {
	r := rule.IsNil[string]()
	s := "hello"
	err := r(&s)
	if err == nil {
		t.Errorf("Error expected for non-nil string pointer")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestIsNilStructNil(t *testing.T) {
	type User struct{ Name string }
	r := rule.IsNil[User]()
	var p *User = nil
	if err := r(p); err != nil {
		t.Errorf("No errors expected for nil struct pointer, got: %v", err)
	}
}

func TestIsNilStructError(t *testing.T) {
	type User struct{ Name string }
	r := rule.IsNil[User]()
	u := User{Name: "Alice"}
	err := r(&u)
	if err == nil {
		t.Errorf("Error expected for non-nil struct pointer")
	} else {
		t.Logf("got error: %v", err)
	}
}

// IsNotNil

func TestIsNotNilInt(t *testing.T) {
	r := rule.IsNotNil[int]()
	v := 99
	if err := r(&v); err != nil {
		t.Errorf("No errors expected for non-nil pointer, got: %v", err)
	}
}

func TestIsNotNilIntError(t *testing.T) {
	r := rule.IsNotNil[int]()
	var p *int = nil
	err := r(p)
	if err == nil {
		t.Errorf("Error expected for nil pointer")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestIsNotNilString(t *testing.T) {
	r := rule.IsNotNil[string]()
	s := "world"
	if err := r(&s); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIsNotNilStringError(t *testing.T) {
	r := rule.IsNotNil[string]()
	var p *string = nil
	err := r(p)
	if err == nil {
		t.Errorf("Error expected for nil string pointer")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestIsNotNilBool(t *testing.T) {
	r := rule.IsNotNil[bool]()
	b := true
	if err := r(&b); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestIsNotNilBoolError(t *testing.T) {
	r := rule.IsNotNil[bool]()
	var p *bool = nil
	err := r(p)
	if err == nil {
		t.Errorf("Error expected for nil bool pointer")
	} else {
		t.Logf("got error: %v", err)
	}
}

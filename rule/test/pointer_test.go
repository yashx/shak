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
	if err := r(&v); err == nil {
		t.Errorf("Error expected for non-nil pointer")
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
	if err := r(&s); err == nil {
		t.Errorf("Error expected for non-nil string pointer")
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
	if err := r(&u); err == nil {
		t.Errorf("Error expected for non-nil struct pointer")
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
	if err := r(p); err == nil {
		t.Errorf("Error expected for nil pointer")
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
	if err := r(p); err == nil {
		t.Errorf("Error expected for nil string pointer")
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
	if err := r(p); err == nil {
		t.Errorf("Error expected for nil bool pointer")
	}
}

package test

import (
	"testing"

	"github.com/yashx/shak/rule"
)

func TestIsTrue(t *testing.T) {
	r := rule.IsTrue[bool]()
	err := r(true)
	if err != nil {
		t.Errorf("No errors are expected, got: %v", err)
	}
}

func TestIsTrueError(t *testing.T) {
	r := rule.IsTrue[bool]()
	err := r(false)
	if err == nil {
		t.Errorf("Errors are expected, got: %v", err)
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestIsFalse(t *testing.T) {
	r := rule.IsFalse[bool]()
	err := r(false)
	if err != nil {
		t.Errorf("No errors are expected, got: %v", err)
	}
}

func TestIsFalseError(t *testing.T) {
	r := rule.IsFalse[bool]()
	err := r(true)
	if err == nil {
		t.Errorf("Errors are expected, got: %v", err)
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestIsTrueWrapper(t *testing.T) {
	type MyBool bool
	var v MyBool = true
	r := rule.IsTrue[MyBool]()
	err := r(v)
	if err != nil {
		t.Errorf("No errors are expected, got: %v", err)
	}
}

func TestIsFalseWrapper(t *testing.T) {
	type MyBool bool
	var v MyBool = false
	r := rule.IsFalse[MyBool]()
	err := r(v)
	if err != nil {
		t.Errorf("No errors are expected, got: %v", err)
	}
}

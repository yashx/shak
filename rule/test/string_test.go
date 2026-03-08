package test

import (
	"regexp"
	"testing"

	"github.com/yashx/shak/rule"
)

// LengthString

func TestLengthStringValid(t *testing.T) {
	r := rule.LengthString[string](5)
	if err := r("hello"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestLengthStringError(t *testing.T) {
	r := rule.LengthString[string](5)
	if err := r("hi"); err == nil {
		t.Errorf("Error expected for wrong byte length")
	}
}

func TestLengthStringEmpty(t *testing.T) {
	r := rule.LengthString[string](0)
	if err := r(""); err != nil {
		t.Errorf("No errors expected for empty string with length 0, got: %v", err)
	}
}

func TestLengthStringCustomType(t *testing.T) {
	type Username string
	r := rule.LengthString[Username](4)
	if err := r(Username("yash")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestLengthStringCustomTypeError(t *testing.T) {
	type Username string
	r := rule.LengthString[Username](3)
	if err := r(Username("yash")); err == nil {
		t.Errorf("Error expected for wrong length custom string type")
	}
}

// RuneLength

func TestRuneLengthASCII(t *testing.T) {
	r := rule.RuneLength[string](5)
	if err := r("hello"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestRuneLengthASCIIError(t *testing.T) {
	r := rule.RuneLength[string](3)
	if err := r("hello"); err == nil {
		t.Errorf("Error expected for wrong rune count")
	}
}

func TestRuneLengthUnicode(t *testing.T) {
	// "héllo" is 5 runes but 6 bytes
	r := rule.RuneLength[string](5)
	if err := r("héllo"); err != nil {
		t.Errorf("No errors expected for 5-rune unicode string, got: %v", err)
	}
}

func TestRuneLengthCustomType(t *testing.T) {
	type Bio string
	r := rule.RuneLength[Bio](3)
	if err := r(Bio("日本語")); err != nil {
		t.Errorf("No errors expected for 3-rune CJK string, got: %v", err)
	}
}

func TestRuneLengthCustomTypeError(t *testing.T) {
	type Bio string
	r := rule.RuneLength[Bio](2)
	if err := r(Bio("日本語")); err == nil {
		t.Errorf("Error expected for wrong rune count on custom type")
	}
}

// MatchString

func TestMatchStringValid(t *testing.T) {
	re := regexp.MustCompile(`^\d{3}-\d{4}$`)
	r := rule.MatchString[string](re)
	if err := r("123-4567"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMatchStringError(t *testing.T) {
	re := regexp.MustCompile(`^\d{3}-\d{4}$`)
	r := rule.MatchString[string](re)
	if err := r("abc"); err == nil {
		t.Errorf("Error expected for non-matching string")
	}
}

func TestMatchStringCustomType(t *testing.T) {
	type Email string
	re := regexp.MustCompile(`^[^@]+@[^@]+$`)
	r := rule.MatchString[Email](re)
	if err := r(Email("user@example.com")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMatchStringCustomTypeError(t *testing.T) {
	type Email string
	re := regexp.MustCompile(`^[^@]+@[^@]+$`)
	r := rule.MatchString[Email](re)
	if err := r(Email("not-an-email")); err == nil {
		t.Errorf("Error expected for invalid email format")
	}
}

// EqualIgnoreCase

func TestEqualIgnoreCaseValid(t *testing.T) {
	r := rule.EqualIgnoreCase[string]("go")
	if err := r("GO"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestEqualIgnoreCaseMixed(t *testing.T) {
	r := rule.EqualIgnoreCase[string]("Hello")
	if err := r("hElLo"); err != nil {
		t.Errorf("No errors expected for mixed case, got: %v", err)
	}
}

func TestEqualIgnoreCaseError(t *testing.T) {
	r := rule.EqualIgnoreCase[string]("go")
	if err := r("rust"); err == nil {
		t.Errorf("Error expected for non-matching string")
	}
}

func TestEqualIgnoreCaseCustomType(t *testing.T) {
	type Lang string
	r := rule.EqualIgnoreCase[Lang]("golang")
	if err := r(Lang("GOLANG")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotEqualIgnoreCase

func TestNotEqualIgnoreCaseValid(t *testing.T) {
	r := rule.NotEqualIgnoreCase[string]("admin")
	if err := r("user"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotEqualIgnoreCaseError(t *testing.T) {
	r := rule.NotEqualIgnoreCase[string]("admin")
	if err := r("ADMIN"); err == nil {
		t.Errorf("Error expected when strings are equal ignoring case")
	}
}

func TestNotEqualIgnoreCaseCaseSensitiveError(t *testing.T) {
	r := rule.NotEqualIgnoreCase[string]("root")
	if err := r("Root"); err == nil {
		t.Errorf("Error expected for case-insensitive match")
	}
}

func TestNotEqualIgnoreCaseCustomType(t *testing.T) {
	type Role string
	r := rule.NotEqualIgnoreCase[Role]("superuser")
	if err := r(Role("guest")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotEqualIgnoreCaseCustomTypeError(t *testing.T) {
	type Role string
	r := rule.NotEqualIgnoreCase[Role]("superuser")
	if err := r(Role("SUPERUSER")); err == nil {
		t.Errorf("Error expected for case-insensitive match on custom type")
	}
}

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
	err := r("hi")
	if err == nil {
		t.Errorf("Error expected for wrong byte length")
	} else {
		t.Logf("got error: %v", err)
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
	err := r(Username("yash"))
	if err == nil {
		t.Errorf("Error expected for wrong length custom string type")
	} else {
		t.Logf("got error: %v", err)
	}
}

// LengthRune

func TestLengthRuneASCII(t *testing.T) {
	r := rule.LengthRune[string](5)
	if err := r("hello"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestLengthRuneASCIIError(t *testing.T) {
	r := rule.LengthRune[string](3)
	err := r("hello")
	if err == nil {
		t.Errorf("Error expected for wrong rune count")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestLengthRuneUnicode(t *testing.T) {
	// "héllo" is 5 runes but 6 bytes
	r := rule.LengthRune[string](5)
	if err := r("héllo"); err != nil {
		t.Errorf("No errors expected for 5-rune unicode string, got: %v", err)
	}
}

func TestLengthRuneCustomType(t *testing.T) {
	type Bio string
	r := rule.LengthRune[Bio](3)
	if err := r(Bio("日本語")); err != nil {
		t.Errorf("No errors expected for 3-rune CJK string, got: %v", err)
	}
}

func TestLengthRuneCustomTypeError(t *testing.T) {
	type Bio string
	r := rule.LengthRune[Bio](2)
	err := r(Bio("日本語"))
	if err == nil {
		t.Errorf("Error expected for wrong rune count on custom type")
	} else {
		t.Logf("got error: %v", err)
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
	err := r("abc")
	if err == nil {
		t.Errorf("Error expected for non-matching string")
	} else {
		t.Logf("got error: %v", err)
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
	err := r(Email("not-an-email"))
	if err == nil {
		t.Errorf("Error expected for invalid email format")
	} else {
		t.Logf("got error: %v", err)
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
	err := r("rust")
	if err == nil {
		t.Errorf("Error expected for non-matching string")
	} else {
		t.Logf("got error: %v", err)
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
	err := r("ADMIN")
	if err == nil {
		t.Errorf("Error expected when strings are equal ignoring case")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotEqualIgnoreCaseCaseSensitiveError(t *testing.T) {
	r := rule.NotEqualIgnoreCase[string]("root")
	err := r("Root")
	if err == nil {
		t.Errorf("Error expected for case-insensitive match")
	} else {
		t.Logf("got error: %v", err)
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
	err := r(Role("SUPERUSER"))
	if err == nil {
		t.Errorf("Error expected for case-insensitive match on custom type")
	} else {
		t.Logf("got error: %v", err)
	}
}

// MinLengthString

func TestMinLengthStringValid(t *testing.T) {
	r := rule.MinLengthString[string](3)
	if err := r("hello"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMinLengthStringAtBound(t *testing.T) {
	r := rule.MinLengthString[string](5)
	if err := r("hello"); err != nil {
		t.Errorf("No errors expected at exact min length, got: %v", err)
	}
}

func TestMinLengthStringError(t *testing.T) {
	r := rule.MinLengthString[string](10)
	err := r("hi")
	if err == nil {
		t.Errorf("Error expected for string below min length")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestMinLengthStringCustomType(t *testing.T) {
	type Username string
	r := rule.MinLengthString[Username](3)
	if err := r(Username("yash")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// MaxLengthString

func TestMaxLengthStringValid(t *testing.T) {
	r := rule.MaxLengthString[string](10)
	if err := r("hello"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMaxLengthStringAtBound(t *testing.T) {
	r := rule.MaxLengthString[string](5)
	if err := r("hello"); err != nil {
		t.Errorf("No errors expected at exact max length, got: %v", err)
	}
}

func TestMaxLengthStringError(t *testing.T) {
	r := rule.MaxLengthString[string](3)
	err := r("hello")
	if err == nil {
		t.Errorf("Error expected for string above max length")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestMaxLengthStringCustomType(t *testing.T) {
	type Username string
	r := rule.MaxLengthString[Username](20)
	if err := r(Username("yash")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// MinLengthRune

func TestMinLengthRuneValid(t *testing.T) {
	r := rule.MinLengthRune[string](3)
	if err := r("héllo"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMinLengthRuneError(t *testing.T) {
	r := rule.MinLengthRune[string](10)
	err := r("hi")
	if err == nil {
		t.Errorf("Error expected for string below min rune length")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestMinLengthRuneUnicode(t *testing.T) {
	r := rule.MinLengthRune[string](3)
	if err := r("日本語"); err != nil {
		t.Errorf("No errors expected for 3-rune CJK string, got: %v", err)
	}
}

func TestMinLengthRuneCustomType(t *testing.T) {
	type Bio string
	r := rule.MinLengthRune[Bio](2)
	if err := r(Bio("日本")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// MaxLengthRune

func TestMaxLengthRuneValid(t *testing.T) {
	r := rule.MaxLengthRune[string](10)
	if err := r("hello"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMaxLengthRuneError(t *testing.T) {
	r := rule.MaxLengthRune[string](2)
	err := r("日本語")
	if err == nil {
		t.Errorf("Error expected for string above max rune length")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestMaxLengthRuneAtBound(t *testing.T) {
	r := rule.MaxLengthRune[string](3)
	if err := r("日本語"); err != nil {
		t.Errorf("No errors expected at exact max rune length, got: %v", err)
	}
}

func TestMaxLengthRuneCustomType(t *testing.T) {
	type Bio string
	r := rule.MaxLengthRune[Bio](5)
	if err := r(Bio("hello")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotMatchString

func TestNotMatchStringValid(t *testing.T) {
	re := regexp.MustCompile(`\d`)
	r := rule.NotMatchString[string](re)
	if err := r("hello"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotMatchStringError(t *testing.T) {
	re := regexp.MustCompile(`\d`)
	r := rule.NotMatchString[string](re)
	err := r("hello123")
	if err == nil {
		t.Errorf("Error expected for string matching pattern")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotMatchStringCustomType(t *testing.T) {
	type Slug string
	re := regexp.MustCompile(`\s`)
	r := rule.NotMatchString[Slug](re)
	if err := r(Slug("no-spaces-here")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotContainsString

func TestNotContainsStringValid(t *testing.T) {
	r := rule.NotContainsString[string]("bad")
	if err := r("hello world"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotContainsStringError(t *testing.T) {
	r := rule.NotContainsString[string]("bad")
	err := r("this is bad")
	if err == nil {
		t.Errorf("Error expected for string containing substr")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotContainsStringCustomType(t *testing.T) {
	type Comment string
	r := rule.NotContainsString[Comment]("spam")
	if err := r(Comment("great post!")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotHasPrefix

func TestNotHasPrefixValid(t *testing.T) {
	r := rule.NotHasPrefix[string]("tmp_")
	if err := r("users"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotHasPrefixError(t *testing.T) {
	r := rule.NotHasPrefix[string]("tmp_")
	err := r("tmp_file")
	if err == nil {
		t.Errorf("Error expected for string with forbidden prefix")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotHasPrefixCustomType(t *testing.T) {
	type TableName string
	r := rule.NotHasPrefix[TableName]("_")
	if err := r(TableName("users")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotHasSuffix

func TestNotHasSuffixValid(t *testing.T) {
	r := rule.NotHasSuffix[string](".tmp")
	if err := r("report.pdf"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotHasSuffixError(t *testing.T) {
	r := rule.NotHasSuffix[string](".tmp")
	err := r("draft.tmp")
	if err == nil {
		t.Errorf("Error expected for string with forbidden suffix")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotHasSuffixCustomType(t *testing.T) {
	type Filename string
	r := rule.NotHasSuffix[Filename](".exe")
	if err := r(Filename("script.sh")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotContainsAnyChar

func TestNotContainsAnyCharValid(t *testing.T) {
	r := rule.NotContainsAnyChar[string](" \t\n")
	if err := r("nospaces"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotContainsAnyCharError(t *testing.T) {
	r := rule.NotContainsAnyChar[string](" \t\n")
	err := r("has space")
	if err == nil {
		t.Errorf("Error expected for string containing forbidden char")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotContainsAnyCharCustomType(t *testing.T) {
	type Password string
	r := rule.NotContainsAnyChar[Password](" ")
	if err := r(Password("s3cr3t!")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotContainsRune

func TestNotContainsRuneValid(t *testing.T) {
	r := rule.NotContainsRune[string]('@')
	if err := r("notanemail"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotContainsRuneError(t *testing.T) {
	r := rule.NotContainsRune[string]('@')
	err := r("user@example.com")
	if err == nil {
		t.Errorf("Error expected for string containing forbidden rune")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotContainsRuneCustomType(t *testing.T) {
	type Tag string
	r := rule.NotContainsRune[Tag](' ')
	if err := r(Tag("no-spaces")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// IsBlank

func TestIsBlankValid(t *testing.T) {
	r := rule.IsBlank[string]()
	if err := r(""); err != nil {
		t.Errorf("No errors expected for empty string, got: %v", err)
	}
}

func TestIsBlankWhitespaceValid(t *testing.T) {
	r := rule.IsBlank[string]()
	if err := r("   "); err != nil {
		t.Errorf("No errors expected for whitespace-only string, got: %v", err)
	}
}

func TestIsBlankError(t *testing.T) {
	r := rule.IsBlank[string]()
	err := r("hello")
	if err == nil {
		t.Errorf("Error expected for non-blank string")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestIsBlankCustomType(t *testing.T) {
	type OptionalField string
	r := rule.IsBlank[OptionalField]()
	if err := r(OptionalField("\t\n")); err != nil {
		t.Errorf("No errors expected for tab/newline-only string, got: %v", err)
	}
}

// ContainsString

func TestContainsStringValid(t *testing.T) {
	r := rule.ContainsString[string]("world")
	if err := r("hello world"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestContainsStringError(t *testing.T) {
	r := rule.ContainsString[string]("world")
	err := r("hello")
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestContainsStringEmpty(t *testing.T) {
	r := rule.ContainsString[string]("")
	if err := r("anything"); err != nil {
		t.Errorf("No errors expected for empty substr, got: %v", err)
	}
}

func TestContainsStringCustomType(t *testing.T) {
	type Bio string
	r := rule.ContainsString[Bio]("go")
	if err := r(Bio("I love golang")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// HasPrefix

func TestHasPrefixValid(t *testing.T) {
	r := rule.HasPrefix[string]("hello")
	if err := r("hello world"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestHasPrefixError(t *testing.T) {
	r := rule.HasPrefix[string]("world")
	err := r("hello world")
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestHasPrefixEmpty(t *testing.T) {
	r := rule.HasPrefix[string]("")
	if err := r("anything"); err != nil {
		t.Errorf("No errors expected for empty prefix, got: %v", err)
	}
}

func TestHasPrefixCustomType(t *testing.T) {
	type Path string
	r := rule.HasPrefix[Path]("/api")
	if err := r(Path("/api/users")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// HasSuffix

func TestHasSuffixValid(t *testing.T) {
	r := rule.HasSuffix[string](".go")
	if err := r("main.go"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestHasSuffixError(t *testing.T) {
	r := rule.HasSuffix[string](".go")
	err := r("main.rs")
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestHasSuffixCustomType(t *testing.T) {
	type Filename string
	r := rule.HasSuffix[Filename](".json")
	if err := r(Filename("config.json")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// ContainsAnyChar

func TestContainsAnyCharValid(t *testing.T) {
	r := rule.ContainsAnyChar[string]("aeiou")
	if err := r("hello"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestContainsAnyCharError(t *testing.T) {
	r := rule.ContainsAnyChar[string]("aeiou")
	err := r("rhythm")
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestContainsAnyCharCustomType(t *testing.T) {
	type Password string
	r := rule.ContainsAnyChar[Password]("!@#$%")
	if err := r(Password("secret!")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// ContainsRune

func TestContainsRuneValid(t *testing.T) {
	r := rule.ContainsRune[string]('!')
	if err := r("hello!"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestContainsRuneError(t *testing.T) {
	r := rule.ContainsRune[string]('!')
	err := r("hello")
	if err == nil {
		t.Errorf("Error expected, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestContainsRuneUnicode(t *testing.T) {
	r := rule.ContainsRune[string]('é')
	if err := r("café"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestContainsRuneCustomType(t *testing.T) {
	type Tag string
	r := rule.ContainsRune[Tag]('-')
	if err := r(Tag("well-formed")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// NotBlank

func TestNotBlankValid(t *testing.T) {
	r := rule.NotBlank[string]()
	if err := r("hello"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotBlankError(t *testing.T) {
	r := rule.NotBlank[string]()
	err := r("")
	if err == nil {
		t.Errorf("Error expected for empty string, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotBlankWhitespaceError(t *testing.T) {
	r := rule.NotBlank[string]()
	err := r("   ")
	if err == nil {
		t.Errorf("Error expected for whitespace-only string, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotBlankTabsAndNewlines(t *testing.T) {
	r := rule.NotBlank[string]()
	err := r("\t\n")
	if err == nil {
		t.Errorf("Error expected for tab/newline-only string, got nil")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotBlankCustomType(t *testing.T) {
	type Name string
	r := rule.NotBlank[Name]()
	if err := r(Name("yash")); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}


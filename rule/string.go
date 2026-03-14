package rule

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/yashx/shak/types"
)

// ContainsString validates that the string contains substr.
func ContainsString[T ~string](substr string) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if strings.Contains(string(value), substr) {
			return nil
		}
		return types.NewValidationError("%q not found in value", substr)
	}
}

// HasPrefix validates that the string has the given prefix.
func HasPrefix[T ~string](prefix string) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if strings.HasPrefix(string(value), prefix) {
			return nil
		}
		return types.NewValidationError("value does not have prefix %q", prefix)
	}
}

// HasSuffix validates that the string has the given suffix.
func HasSuffix[T ~string](suffix string) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if strings.HasSuffix(string(value), suffix) {
			return nil
		}
		return types.NewValidationError("value does not have suffix %q", suffix)
	}
}

// ContainsAnyChar validates that the string contains any Unicode code point in chars.
func ContainsAnyChar[T ~string](chars string) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if strings.ContainsAny(string(value), chars) {
			return nil
		}
		return types.NewValidationError("value does not contain any of %q", chars)
	}
}

// ContainsRune validates that the string contains the given rune.
func ContainsRune[T ~string](r rune) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if strings.ContainsRune(string(value), r) {
			return nil
		}
		return types.NewValidationError("value does not contain %q", r)
	}
}

// NotBlank validates that the string is not empty or whitespace-only.
func NotBlank[T ~string]() types.Rule[T] {
	return func(value T) *types.ValidationError {
		if strings.TrimSpace(string(value)) != "" {
			return nil
		}
		return types.NewValidationError("value must not be blank")
	}
}

// LengthString validates that the string has exactly l bytes. Use LengthRune for character count.
func LengthString[T ~string](l int) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if len(value) == l {
			return nil
		}
		return types.NewValidationError("length is %d bytes, expected %d", len(value), l)
	}
}

// LengthRune validates that the string contains exactly l Unicode code points (runes).
// Note: a single visible character may span multiple code points. For example,
// the family emoji 👨‍👩‍👧 is composed of several code points joined by zero-width
// joiners, so LengthRune counts it as more than 1. If you need to validate
// display length (e.g. for usernames shown in a UI), use a grapheme-cluster-aware
// library such as github.com/rivo/uniseg instead.
func LengthRune[T ~string](l int) types.Rule[T] {
	return func(value T) *types.ValidationError {
		count := utf8.RuneCountInString(string(value))
		if count == l {
			return nil
		}
		return types.NewValidationError("length is %d characters, expected %d", count, l)
	}
}

// MatchString validates that the string matches the provided regular expression.
func MatchString[T ~string](re *regexp.Regexp) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if re.MatchString(string(value)) {
			return nil
		}
		return types.NewValidationError("value does not match the required format")
	}
}

// EqualIgnoreCase validates that the string equals eq, ignoring case.
func EqualIgnoreCase[T ~string](eq string) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if strings.EqualFold(string(value), eq) {
			return nil
		}
		return types.NewValidationError("ignoring case, %v is not equal to %v", value, eq)
	}
}

// NotEqualIgnoreCase validates that the string does not equal eq, ignoring case.
func NotEqualIgnoreCase[T ~string](eq string) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if !strings.EqualFold(string(value), eq) {
			return nil
		}
		return types.NewValidationError("ignoring case, %v must not equal %v", value, eq)
	}
}

// MinLengthString validates that the string has at least l bytes.
func MinLengthString[T ~string](l int) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if len(value) >= l {
			return nil
		}
		return types.NewValidationError("length is %d bytes, expected at least %d", len(value), l)
	}
}

// MaxLengthString validates that the string has at most l bytes.
func MaxLengthString[T ~string](l int) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if len(value) <= l {
			return nil
		}
		return types.NewValidationError("length is %d bytes, expected at most %d", len(value), l)
	}
}

// MinLengthRune validates that the string contains at least l Unicode code points.
func MinLengthRune[T ~string](l int) types.Rule[T] {
	return func(value T) *types.ValidationError {
		count := utf8.RuneCountInString(string(value))
		if count >= l {
			return nil
		}
		return types.NewValidationError("length is %d characters, expected at least %d", count, l)
	}
}

// MaxLengthRune validates that the string contains at most l Unicode code points.
func MaxLengthRune[T ~string](l int) types.Rule[T] {
	return func(value T) *types.ValidationError {
		count := utf8.RuneCountInString(string(value))
		if count <= l {
			return nil
		}
		return types.NewValidationError("length is %d characters, expected at most %d", count, l)
	}
}

// NotMatchString validates that the string does not match the provided regular expression.
func NotMatchString[T ~string](re *regexp.Regexp) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if !re.MatchString(string(value)) {
			return nil
		}
		return types.NewValidationError("value must not match the pattern %q", re.String())
	}
}

// NotContainsString validates that the string does not contain substr.
func NotContainsString[T ~string](substr string) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if !strings.Contains(string(value), substr) {
			return nil
		}
		return types.NewValidationError("%q must not appear in value", substr)
	}
}

// NotHasPrefix validates that the string does not have the given prefix.
func NotHasPrefix[T ~string](prefix string) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if !strings.HasPrefix(string(value), prefix) {
			return nil
		}
		return types.NewValidationError("value must not have prefix %q", prefix)
	}
}

// NotHasSuffix validates that the string does not have the given suffix.
func NotHasSuffix[T ~string](suffix string) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if !strings.HasSuffix(string(value), suffix) {
			return nil
		}
		return types.NewValidationError("value must not have suffix %q", suffix)
	}
}

// NotContainsAnyChar validates that the string contains none of the Unicode code points in chars.
func NotContainsAnyChar[T ~string](chars string) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if !strings.ContainsAny(string(value), chars) {
			return nil
		}
		return types.NewValidationError("value must not contain any of %q", chars)
	}
}

// NotContainsRune validates that the string does not contain the given rune.
func NotContainsRune[T ~string](r rune) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if !strings.ContainsRune(string(value), r) {
			return nil
		}
		return types.NewValidationError("value must not contain %q", r)
	}
}

// IsBlank validates that the string is empty or whitespace-only.
func IsBlank[T ~string]() types.Rule[T] {
	return func(value T) *types.ValidationError {
		if strings.TrimSpace(string(value)) == "" {
			return nil
		}
		return types.NewValidationError("value must be blank")
	}
}

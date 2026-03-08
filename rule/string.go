package rule

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/yashx/shak/types"
)

// LengthString validates that the string has exactly l bytes. Use RuneLength for character count.
func LengthString[T ~string](l int) types.Rule[T] {
	return func(value T) *types.ValidationError {
		if len(value) == l {
			return nil
		}
		return types.NewValidationError("length is %d bytes, expected %d", len(value), l)
	}
}

// RuneLength validates that the string contains exactly l Unicode code points (runes).
// Note: a single visible character may span multiple code points. For example,
// the family emoji 👨‍👩‍👧 is composed of several code points joined by zero-width
// joiners, so RuneLength counts it as more than 1. If you need to validate
// display length (e.g. for usernames shown in a UI), use a grapheme-cluster-aware
// library such as github.com/rivo/uniseg instead.
func RuneLength[T ~string](l int) types.Rule[T] {
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

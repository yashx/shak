package types

import (
	"fmt"
	"strings"
)

// ValidationError holds a validation failure with the path to the offending
// field and a human-readable message.
type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (er ValidationError) Error() string {
	if er.Path == "" || er.Path == "." {
		return er.Message
	}
	return er.Path + ": " + er.Message
}

// NewValidationError creates a ValidationError with path "." and a formatted message.
func NewValidationError(messageFormat string, a ...any) *ValidationError {
	return &ValidationError{
		Path:    ".",
		Message: fmt.Sprintf(messageFormat, a...),
	}
}

// NewValidationErrorWithPath creates a ValidationError with an explicit path and message.
func NewValidationErrorWithPath(path, message string) *ValidationError {
	return &ValidationError{
		Path:    path,
		Message: message,
	}
}

// NewValidationErrorWithParentName prefixes childErr's path with parentName.
// If parentName is "" or ".", the child path is returned unchanged.
func NewValidationErrorWithParentName(parentName string, childErr *ValidationError) *ValidationError {
	var p string
	if parentName == "" || parentName == "." {
		p = childErr.Path
	} else if strings.HasPrefix(childErr.Path, "[") {
		p = fmt.Sprintf("%s%s", parentName, childErr.Path)
	} else {
		p = fmt.Sprintf("%s.%s", parentName, childErr.Path)
	}
	p = strings.TrimRight(p, ".")
	return &ValidationError{
		Path:    p,
		Message: childErr.Message,
	}
}

// Rule is a validation function for a value of type T.
// It returns nil on success and a ValidationError on failure.
type Rule[T any] func(T) *ValidationError

// WithMessage returns a new rule that runs r but replaces its error message,
// preserving the original error path.
func (r Rule[T]) WithMessage(format string, a ...any) Rule[T] {
	return func(value T) *ValidationError {
		err := r(value)
		if err == nil {
			return nil
		}
		return &ValidationError{Path: err.Path, Message: fmt.Sprintf(format, a...)}
	}
}

// Else returns r if r is non-nil, otherwise returns rule.
// Useful as a fallback when If returns nil.
func (r Rule[T]) Else(rule Rule[T]) Rule[T] {
	if r != nil {
		return r
	}

	return rule
}

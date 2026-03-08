package types

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func (er ValidationError) Error() string {
	return er.Message
}

func NewValidationError(messageFormat string, a ...any) *ValidationError {
	return &ValidationError{
		Path:    ".",
		Message: fmt.Sprintf(messageFormat, a...),
	}
}

func NewValidationErrorWithPath(path, message string) *ValidationError {
	return &ValidationError{
		Path:    path,
		Message: message,
	}
}

func NewValidationErrorWithParentName(parentName string, childErr *ValidationError) *ValidationError {
	var p string
	if strings.HasPrefix(childErr.Path, "[") {
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

type Rule[T any] func(T) *ValidationError

func (r Rule[T]) WithMessage(format string, a ...any) Rule[T] {
	return func(value T) *ValidationError {
		err := r(value)
		if err == nil {
			return nil
		}
		return NewValidationError(format, a...)
	}
}

func (r Rule[T]) Else(rule Rule[T]) Rule[T] {
	if r != nil {
		return r
	}

	return rule
}

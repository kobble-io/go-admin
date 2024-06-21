package utils

import "fmt"

type ErrorBase struct {
	Name    string
	Message string
	Cause   error
}

func NewErrorBase(name, message string, cause error) *ErrorBase {
	return &ErrorBase{
		Name:    name,
		Message: message,
		Cause:   cause,
	}
}

func (e *ErrorBase) Error() string {
	return fmt.Sprintf("%s: %s", e.Name, e.Message)
}

package commonerrs

import "fmt"

type InvalidInputError struct {
	Message string
}

func (e InvalidInputError) Error() string {
	return fmt.Sprintf("invalid input: %s", e.Message)
}

func NewInvalidInputError(message string) error {
	return InvalidInputError{Message: message}
}

func NewInvalidInputErrorf(format string, a ...interface{}) error {
	return InvalidInputError{Message: fmt.Sprintf(format, a...)}
}

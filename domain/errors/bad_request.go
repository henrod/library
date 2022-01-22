package errors

import "fmt"

type BadRequestError struct {
	InvalidField string
	Details      string
}

func (n BadRequestError) Error() string {
	return fmt.Sprintf("field %s is invalid: %s", n.InvalidField, n.Details)
}

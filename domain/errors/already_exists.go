package errors

import "fmt"

type AlreadyExistsError struct {
	Details string
}

func (n AlreadyExistsError) Error() string {
	return fmt.Sprintf("resource already exists: %s", n.Details)
}

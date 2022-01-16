package errors

import "fmt"

type NotFoundError struct {
	Details string
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", n.Details)
}

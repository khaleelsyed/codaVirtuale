package storage

import (
	"errors"
	"fmt"
)

var ErrNotImplemented = errors.New("not implemented")
var ErrNotFound = errors.New("not found")
var ErrFailedToDelete = errors.New("failed to delete")

func errAffectedMultipleRows(operation string) error {
	return fmt.Errorf("multiple rows were affected during %s", operation)
}

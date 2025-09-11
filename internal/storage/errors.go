package storage

import (
	"errors"
	"fmt"
)

var ErrNotImplemented = errors.New("not implemented")
var ErrNotFound = errors.New("not found")
var ErrFailedToDelete = errors.New("failed to delete")
var errDeskNull = errors.New("sql: Scan error on column index 3, name \"desk_id\": converting NULL to int is unsupported")

func errAffectedMultipleRows(operation string) error {
	return fmt.Errorf("multiple rows were affected during %s", operation)
}

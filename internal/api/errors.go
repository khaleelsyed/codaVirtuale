package api

import (
	"errors"
	"fmt"
)

const pqUniqueConstraintViolation = "pq: duplicate key value violates unique constraint"
const pqForeignKeyConstraintViolation = "violates foreign key constraint"

var errBadRequestBody = errors.New("bad request body")

func badValidationString(object_type string) string {
	return fmt.Sprintf("error validating %s", object_type)
}

type apiError struct {
	s string
}

func (e apiError) Error() string {
	return e.s
}

func APIError(s string) error {
	return apiError{s: s}
}

package api

import (
	"errors"
	"fmt"
)

var errBadRequestBody = errors.New("bad request body")
var errNotFound = errors.New("not found")

func badValidationString(object_type string) string {
	return fmt.Sprintf("error validating %s", object_type)
}

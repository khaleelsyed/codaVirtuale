package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func writeJSON(w http.ResponseWriter, status int, v any, logger *zap.SugaredLogger) error {
	w.Header().Set("Content-Type", "application/json")

	if isErrorValue(v) {
		v, status = handleErrors(v, status, logger)
	}

	w.WriteHeader(status)

	if v != nil {
		return json.NewEncoder(w).Encode(v)
	}
	return nil
}

func handleErrors(v any, status int, logger *zap.SugaredLogger) (any, int) {
	if err, found := v.(error); found {
		validateErrStatus(status, err, logger)

		return err.Error(), status
	}

	if errs, found := v.([]error); found {

		validateErrStatus(status, errs, logger)

		errors := make([]string, len(errs))
		for i, err := range errs {
			errors[i] = err.Error()
		}

		return map[string][]string{
			"errors": errors,
		}, status
	}

	logger.Warnw("non error value passed to handleErrors", "value", v, "original status", status)
	return v, status
}

func isErrorValue(v any) bool {
	var found bool

	if _, found = v.([]error); found {
		return true
	} else if _, found = v.(error); found {
		return true
	}

	return false
}

func validateErrStatus(status int, v any, logger *zap.SugaredLogger) int {
	if status < 400 {
		logger.Warnw("unhandled error passed with non-error status code", "original status code", status, "error value", v)
		return http.StatusInternalServerError
	}

	return status
}

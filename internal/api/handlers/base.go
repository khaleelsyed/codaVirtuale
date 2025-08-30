package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandler(f handlerFunc, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			logger.Error("Unhandled error", zap.Error(err))
			writeJSON(w, http.StatusInternalServerError, err, logger)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any, logger *zap.Logger) error {
	w.Header().Set("Content-Type", "application/json")

	if _, errFound := v.(error); errFound {
		if status < 400 {
			status = http.StatusInternalServerError
		}
	}

	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

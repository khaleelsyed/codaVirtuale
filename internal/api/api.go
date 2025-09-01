package api

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

type APIServer struct {
	listenAddress      string
	storage            Storage
	logger             *zap.Logger
	rollingQueueNumber int
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	staffRouter := router.PathPrefix("/internal").Subrouter()
	s.addStaffRoutes(staffRouter)

	s.logger.Info("Listening to requests", zap.String("listenAddress", s.listenAddress))
	if err := http.ListenAndServe(s.listenAddress, router); err != nil {
		s.logger.Error("Failed to run ListenAndServe", zap.Error(err))
	}
}

func makeHTTPHandler(f handlerFunc, allowedMethods []string, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !slices.Contains(allowedMethods, r.Method) {
			writeJSON(w, http.StatusMethodNotAllowed, "method not allowed", logger)
			return
		}

		if err := f(w, r); err != nil {
			logger.Error("Unhandled error", zap.Error(err))
			writeJSON(w, http.StatusInternalServerError, err, logger)
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any, logger *zap.Logger) error {
	w.Header().Set("Content-Type", "application/json")

	if err, errFound := v.(error); errFound {
		if status < 400 {
			logger.Error("unhandled error passed with normal status code", zap.Int("original status code", status), zap.Error(err))
			status = http.StatusInternalServerError
		}
	}

	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func NewAPIServer(listenAddress string, storage Storage, logger *zap.Logger) *APIServer {
	return &APIServer{
		listenAddress:      listenAddress,
		storage:            storage,
		logger:             logger,
		rollingQueueNumber: 1,
	}
}

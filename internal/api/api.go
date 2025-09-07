package api

import (
	"net/http"
	"slices"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

type APIServer struct {
	listenAddress      string
	storage            Storage
	logger             *zap.SugaredLogger
	rollingQueueNumber int
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	staffRouter := router.PathPrefix("/internal").Subrouter()
	s.addStaffRoutes(staffRouter)

	ticketRouter := router.PathPrefix("/ticket").Subrouter()
	s.addTicketRoutes(ticketRouter)

	categoryRouter := router.PathPrefix("/category").Subrouter()
	s.addCategoryRoutes(categoryRouter)

	deskRouter := router.PathPrefix("/desk").Subrouter()
	s.addDeskRoutes(deskRouter)

	s.logger.Infow("Listening to requests", "listenAddress", s.listenAddress)
	if err := http.ListenAndServe(s.listenAddress, router); err != nil {
		s.logger.Errorw("Failed to run ListenAndServe", "error", err)
	}
}

func makeHTTPHandler(f handlerFunc, allowedMethods []string, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !slices.Contains(allowedMethods, r.Method) {
			writeJSON(w, http.StatusMethodNotAllowed, "method not allowed", logger)
			return
		}

		if err := f(w, r); err != nil {
			logger.Errorw("Unhandled error", "error", err)
			writeJSON(w, http.StatusInternalServerError, err, logger)
			return
		}
	}
}

func NewAPIServer(listenAddress string, storage Storage, logger *zap.SugaredLogger) *APIServer {
	return &APIServer{
		listenAddress:      listenAddress,
		storage:            storage,
		logger:             logger,
		rollingQueueNumber: 1,
	}
}

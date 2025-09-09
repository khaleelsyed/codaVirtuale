package api

import (
	"net/http"
	"slices"

	"github.com/gorilla/mux"
	"github.com/khaleelsyed/codaVirtuale/internal/types"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

type APIServer struct {
	listenAddress      string
	storage            Storage
	logger             *types.SugarWithTrace
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

func makeHTTPHandler(f handlerFunc, allowedMethods []string, logger *types.SugarWithTrace) http.HandlerFunc {
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

func NewAPIServer(listenAddress string, storage Storage, logger *types.SugarWithTrace) *APIServer {
	return &APIServer{
		listenAddress:      listenAddress,
		storage:            storage,
		logger:             logger,
		rollingQueueNumber: 1,
	}
}

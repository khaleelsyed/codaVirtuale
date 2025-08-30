package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type APIServer struct {
	listenAddress string
	storage       Storage
	logger        *zap.Logger
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	s.logger.Info("Listening to requests", zap.String("listenAddress", s.listenAddress))
	if err := http.ListenAndServe(s.listenAddress, router); err != nil {
		s.logger.Error("Failed to run ListenAndServe", zap.Error(err))
	}
}

func NewAPIServer(listenAddress string, storage Storage, logger *zap.Logger) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		storage:       storage,
		logger:        logger,
	}
}

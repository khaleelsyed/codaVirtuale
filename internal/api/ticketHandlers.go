package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (s *APIServer) addTicketRoutes(router *mux.Router) {
	router.HandleFunc("/{id}", makeHTTPHandler(s.handleTicket, []string{http.MethodGet, http.MethodDelete}, s.logger))
	router.HandleFunc("", makeHTTPHandler(s.createTicket, []string{http.MethodPost}, s.logger))
}

func (s *APIServer) getTicket(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	ticketID, err := strconv.Atoi(idStr)
	if err != nil {
		errBody := "bad ID"
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	ticket, err := s.storage.GetTicket(ticketID)
	if err != nil {
		errBody := "error retrieving ticket"
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusInternalServerError, errBody, s.logger)
	}

	return writeJSON(w, http.StatusOK, ticket, s.logger)
}

func (s *APIServer) deleteTicket(w http.ResponseWriter, r *http.Request) error {
	var err error

	idStr := mux.Vars(r)["id"]
	ticketID, err := strconv.Atoi(idStr)
	if err != nil {
		errBody := "bad ID"
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	if err = s.storage.DeleteTicket(ticketID); err != nil {
		errBody := "error deleting ticket"
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusInternalServerError, errBody, s.logger)
	}

	return writeJSON(w, http.StatusNoContent, nil, s.logger)
}

func (s *APIServer) createTicket(w http.ResponseWriter, r *http.Request) error {
	var err error

	var requestBody struct {
		CategoryID int `json:"category_id"`
	}

	if err = json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		s.logger.Error("bad request body", zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBadRequestBody, s.logger)
	}

	if _, err = s.storage.GetCategory(requestBody.CategoryID); err != nil {
		errBody := badValidationString("category")
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errors.New(errBody), s.logger)
	}

	s.logger.Debug("", zap.Int("category ID", requestBody.CategoryID))

	ticket, err := s.storage.CreateTicket(requestBody.CategoryID)
	if err != nil {
		errBody := "error creating ticket"
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusInternalServerError, errBody, s.logger)
	}

	return writeJSON(w, http.StatusCreated, ticket, s.logger)
}

func (s *APIServer) handleTicket(w http.ResponseWriter, r *http.Request) error {
	s.logger.Debug("", zap.String("method", r.Method))
	switch r.Method {
	case http.MethodGet:
		return s.getTicket(w, r)
	case http.MethodDelete:
		return s.deleteTicket(w, r)
	default:
		s.logger.Error("unhandled method", zap.String("method", r.Method))
		return writeJSON(w, http.StatusInternalServerError, fmt.Errorf("unhandled method %s", r.Method), s.logger)
	}
}

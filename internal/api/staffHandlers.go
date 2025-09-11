package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) addStaffRoutes(router *mux.Router) {
	router.HandleFunc("/next", makeHTTPHandler(s.handleNext, []string{http.MethodPut, http.MethodGet}, s.logger))
	router.HandleFunc("/queue", makeHTTPHandler(s.getQueue, []string{http.MethodGet}, s.logger))
}

func (s *APIServer) putNextTicket(w http.ResponseWriter, r *http.Request) error {
	var requestBody struct {
		DeskID int `json:"desk_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return writeJSON(w, http.StatusBadRequest, errBadRequestBody, s.logger)
	}

	if _, err := s.storage.GetDesk(requestBody.DeskID); err != nil {
		errBody := badValidationString("desk")
		return writeJSON(w, http.StatusBadRequest, errors.New(errBody), s.logger)
	}

	nextTicket, err := s.storage.CallNextTicket(requestBody.DeskID)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, err, s.logger)
	}

	return writeJSON(w, http.StatusOK, nextTicket, s.logger)
}

func (s *APIServer) getNext(w http.ResponseWriter, r *http.Request) error {
	categoryIDStr := r.URL.Query().Get("category_id")
	if categoryIDStr == "" {
		errBody := "no category_id given as query parameters"
		return writeJSON(w, http.StatusBadRequest, errors.New(errBody), s.logger)
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		errBody := "bad category ID"
		return writeJSON(w, http.StatusBadRequest, errors.New(errBody), s.logger)
	}

	if _, err := s.storage.GetCategory(categoryID); err != nil {
		errBody := badValidationString("category")
		return writeJSON(w, http.StatusBadRequest, errors.New(errBody), s.logger)
	}

	nextTicket, err := s.storage.SeeNext(categoryID)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, err, s.logger)
	}

	return writeJSON(w, http.StatusOK, nextTicket, s.logger)
}

func (s *APIServer) getQueue(w http.ResponseWriter, r *http.Request) error {
	ticketIDs, err := s.storage.SeeQueue()
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, err, s.logger)
	}
	return writeJSON(w, http.StatusOK, ticketIDs, s.logger)
}

func (s *APIServer) handleNext(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.getNext(w, r)
	case http.MethodPut:
		return s.putNextTicket(w, r)
	default:
		s.logger.Errorw("unhandled method", "method", r.Method)
		return writeJSON(w, http.StatusBadRequest, fmt.Errorf("unhandled method %s", r.Method), s.logger)
	}
}

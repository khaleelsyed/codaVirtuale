package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (s *APIServer) addStaffRoutes(router *mux.Router) {
	router.HandleFunc("/next", makeHTTPHandler(s.PutNextTicket, []string{http.MethodPut}, s.logger))
	router.HandleFunc("/last", makeHTTPHandler(s.GetLastCalled, []string{http.MethodGet}, s.logger))
	router.HandleFunc("/peek", makeHTTPHandler(s.GetNext, []string{http.MethodGet}, s.logger))
	router.HandleFunc("/queue", makeHTTPHandler(s.GetQueue, []string{http.MethodGet}, s.logger))
}

func (s *APIServer) PutNextTicket(w http.ResponseWriter, r *http.Request) error {
	var requestBody struct {
		DeskID int `json:"desk_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		s.logger.Error("bad desk ID passed", zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBadRequestBody, s.logger)
	}

	if _, err := s.storage.GetDesk(requestBody.DeskID); err != nil {
		errBody := badValidationString("desk")
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errors.New(errBody), s.logger)
	}

	nextTicket, err := s.storage.CallNextTicket(requestBody.DeskID)
	if err != nil {
		s.logger.Error("error retrieving next ticket", zap.Error(err))
		return writeJSON(w, http.StatusInternalServerError, err, s.logger)
	}

	return writeJSON(w, http.StatusOK, nextTicket, s.logger)
}

func (s *APIServer) GetLastCalled(w http.ResponseWriter, r *http.Request) error {
	var err error

	var requestBody struct {
		CategoryID int `json:"category_id"`
		Positions  int `json:"positions"`
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

	tickets, err := s.storage.LastCalled(requestBody.CategoryID, requestBody.Positions)
	if err != nil {
		s.logger.Error("error retreiving last called tickets", zap.Error(err))
		return writeJSON(w, http.StatusInternalServerError, err, s.logger)
	}

	ticketIDs := make([]int, requestBody.Positions)
	for i, ticket := range tickets {
		ticketIDs[i] = ticket.QueueNumber
	}

	return writeJSON(w, http.StatusOK, ticketIDs, s.logger)
}

func (s *APIServer) GetNext(w http.ResponseWriter, r *http.Request) error {
	var requestBody struct {
		CategoryID int `json:"category_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		s.logger.Error("bad request body", zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBadRequestBody, s.logger)
	}

	if _, err := s.storage.GetCategory(requestBody.CategoryID); err != nil {
		errBody := badValidationString("category")
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errors.New(errBody), s.logger)
	}

	nextTicket, err := s.storage.SeeNext(requestBody.CategoryID)
	if err != nil {
		s.logger.Error("error retreiving next ticket", zap.Error(err))
		return writeJSON(w, http.StatusInternalServerError, err, s.logger)
	}

	return writeJSON(w, http.StatusOK, nextTicket, s.logger)
}

func (s *APIServer) GetQueue(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusNotImplemented, nil, s.logger)
}

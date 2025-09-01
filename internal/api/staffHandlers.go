package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (s *APIServer) addStaffRoutes(router *mux.Router) {
	router.HandleFunc("/next", makeHTTPHandler(s.CallNextTicket, []string{http.MethodPut}, s.logger))
	router.HandleFunc("/last", makeHTTPHandler(s.LastCalled, []string{http.MethodGet}, s.logger))
}

func (s *APIServer) CallNextTicket(w http.ResponseWriter, r *http.Request) error {
	var requestBody struct {
		DeskID int `json:"desk_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		s.logger.Error("bad desk ID passed", zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBadID, s.logger)
	}

	desk, err := s.storage.GetDesk(requestBody.DeskID)
	if err != nil {
		s.logger.Error("error retrieving desk", zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBadID, s.logger)
	}

	nextTicket, err := s.storage.CallNextTicket(desk)
	if err != nil {
		s.logger.Error("error retrieving next ticket", zap.Error(err))
		return writeJSON(w, http.StatusInternalServerError, err, s.logger)
	}

	return writeJSON(w, http.StatusOK, nextTicket, s.logger)
}

func (s *APIServer) LastCalled(w http.ResponseWriter, r *http.Request) error {
	var err error

	var requestBody struct {
		CategoryID int `json:"category_id"`
		Positions  int `json:"positions"`
	}

	if err = json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		s.logger.Error("bad request body", zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, err, s.logger)
	}

	if _, err = s.storage.GetCategory(requestBody.CategoryID); err != nil {
		s.logger.Error("error validating category", zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, err, s.logger)
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

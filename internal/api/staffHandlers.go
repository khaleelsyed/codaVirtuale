package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (s *APIServer) addStaffRoutes(router *mux.Router) {
	router.HandleFunc("/next", makeHTTPHandler(s.CallNextTicket, []string{http.MethodPut}, s.logger))
}

func (s *APIServer) CallNextTicket(w http.ResponseWriter, r *http.Request) error {
	deskID, err := getID(r)
	if err != nil {
		s.logger.Error("bad desk ID passed", zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBadID, s.logger)
	}

	desk, err := s.storage.GetDesk(deskID)
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

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, err
	}

	return id, nil
}

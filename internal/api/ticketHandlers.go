package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/khaleelsyed/codaVirtuale/internal/types"
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
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	ticket, err := s.storage.GetTicket(ticketID)
	if err != nil {
		errBody := badValidationString("category")
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
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	if err = s.storage.DeleteTicket(ticketID); err != nil {
		errBody := "error deleting ticket"
		return writeJSON(w, http.StatusInternalServerError, errBody, s.logger)
	}

	return writeJSON(w, http.StatusNoContent, nil, s.logger)
}

func (s *APIServer) createTicket(w http.ResponseWriter, r *http.Request) error {
	var err error

	var requestBody struct {
		CategoryID int `json:"category_id"`
	}

	randomString := func() string {
		b := make([]byte, 6)
		rand.Read(b)
		return hex.EncodeToString(b)
	}

	if err = json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return writeJSON(w, http.StatusBadRequest, errBadRequestBody, s.logger)
	}

	if requestBody.CategoryID == 0 {
		return writeJSON(w, http.StatusBadRequest, apiError{"bad category ID"}, s.logger)
	}

	_, err = s.storage.GetCategory(requestBody.CategoryID)
	if err != nil {
		if err == types.ErrnotFound {
			return writeJSON(w, http.StatusNotFound, "category not found", s.logger)
		}
		s.logger.Tracew("failed to validate category", "category_id", requestBody.CategoryID, "error", err)
		return writeJSON(w, http.StatusBadRequest, badValidationString("category"), s.logger)
	}

	var ticket types.Ticket

	for range 3 {
		ticket, err = s.storage.CreateTicket(types.TicketCreate{CategoryID: requestBody.CategoryID, SubURL: randomString()})

		if err != nil {
			if strings.HasPrefix(err.Error(), pqUniqueConstraintViolation) {
				s.logger.Tracew("Failed to create a unique ticket url, retrying", "error", err, "sub_url", ticket.SubURL, "category_id", ticket.CategoryID)
				continue
			}
			return writeJSON(w, http.StatusInternalServerError, "error creating ticket", s.logger)
		}

		return writeJSON(w, http.StatusCreated, ticket, s.logger)
	}

	s.logger.Warn("Retry threshold has been reached for generating ticket SubURL")
	return writeJSON(w, http.StatusInternalServerError, errors.New("error creating ticket"), s.logger)
}

func (s *APIServer) handleTicket(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.getTicket(w, r)
	case http.MethodDelete:
		return s.deleteTicket(w, r)
	default:
		s.logger.Errorw("unhandled method", "method", r.Method)
		return writeJSON(w, http.StatusInternalServerError, fmt.Errorf("unhandled method %s", r.Method), s.logger)
	}
}

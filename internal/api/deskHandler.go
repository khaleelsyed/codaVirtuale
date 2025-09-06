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

func (s *APIServer) addDeskRoutes(router *mux.Router) {
	router.HandleFunc("/{id}", makeHTTPHandler(s.handleDesk, []string{http.MethodGet, http.MethodPut, http.MethodDelete}, s.logger))
	router.HandleFunc("", makeHTTPHandler(s.createDesk, []string{http.MethodPost}, s.logger))
}

func (s *APIServer) getDesk(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	deskID, err := strconv.Atoi(idStr)
	if err != nil {
		errBody := "bad ID"
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	desk, err := s.storage.GetDesk(deskID)
	if err != nil {
		errBody := badValidationString("desk")
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	return writeJSON(w, http.StatusOK, desk, s.logger)
}

func (s *APIServer) putDesk(w http.ResponseWriter, r *http.Request) error {
	var err error

	var requestBody struct {
		CategoryID int    `json:"category_id"`
		Label      string `json:"label"`
	}

	idStr := mux.Vars(r)["id"]
	deskID, err := strconv.Atoi(idStr)
	if err != nil {
		errBody := "bad ID"
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	_, err = s.storage.GetDesk(deskID)
	if err != nil {
		errBody := badValidationString("category")
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	if err = json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		s.logger.Error("bad request body", zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBadRequestBody, s.logger)
	}

	s.logger.Debug("", zap.Any("request body", requestBody))

	// DEV
	// category, err := s.storage.UpdateCategory(categoryID, requestBody.Name)
	// if err != nil {
	// 	errBody := badValidationString("category")
	// 	s.logger.Error(errBody, zap.Error(err))
	// 	return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	// }

	// return writeJSON(w, http.StatusOK, category, s.logger)

	return nil
}

func (s *APIServer) deleteDesk(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		errBody := "bad ID"
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	if err := s.storage.DeleteCategory(categoryID); err != nil {
		errBody := badValidationString("category")
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	return writeJSON(w, http.StatusNoContent, nil, s.logger)
}

func (s *APIServer) createDesk(w http.ResponseWriter, r *http.Request) error {
	var err error

	var requestBody struct {
		Name string
	}

	if err = json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		s.logger.Error("bad request body", zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBadRequestBody, s.logger)
	}

	category, err := s.storage.CreateCategory(requestBody.Name)
	if err != nil {
		errBody := "error creating category"
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusInternalServerError, errors.New(errBody), s.logger)
	}

	return writeJSON(w, http.StatusCreated, category, s.logger)
}

func (s *APIServer) handleDesk(w http.ResponseWriter, r *http.Request) error {
	s.logger.Debug("", zap.String("method", r.Method))
	switch r.Method {
	case http.MethodGet:
		return s.getDesk(w, r)
	case http.MethodPut:
		return s.putDesk(w, r)
	case http.MethodDelete:
		return s.deleteDesk(w, r)
	default:
		s.logger.Error("unhandled method", zap.String("method", r.Method))
		return writeJSON(w, http.StatusInternalServerError, fmt.Errorf("unhandled method %s", r.Method), s.logger)
	}
}

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

func (s *APIServer) addCategoryRoutes(router *mux.Router) {
	router.HandleFunc("/{id}", makeHTTPHandler(s.handleCategory, []string{http.MethodGet, http.MethodPut, http.MethodDelete}, s.logger))
	router.HandleFunc("", makeHTTPHandler(s.createCategory, []string{http.MethodPost}, s.logger))
}

func (s *APIServer) getCategory(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		errBody := "bad ID"
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	category, err := s.storage.GetCategory(categoryID)
	if err != nil {
		errBody := badValidationString("category")
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	return writeJSON(w, http.StatusOK, category, s.logger)
}

func (s *APIServer) putCategory(w http.ResponseWriter, r *http.Request) error {
	var err error

	var requestBody struct {
		Name string
	}

	idStr := mux.Vars(r)["id"]
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		errBody := "bad ID"
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	if err = json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		s.logger.Error("bad request body", zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBadRequestBody, s.logger)
	}

	category, err := s.storage.UpdateCategory(categoryID, requestBody.Name)
	if err != nil {
		errBody := badValidationString("category")
		s.logger.Error(errBody, zap.Error(err))
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	return writeJSON(w, http.StatusOK, category, s.logger)
}

func (s *APIServer) deleteCategory(w http.ResponseWriter, r *http.Request) error {
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

func (s *APIServer) createCategory(w http.ResponseWriter, r *http.Request) error {
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

func (s *APIServer) handleCategory(w http.ResponseWriter, r *http.Request) error {
	s.logger.Debug("", zap.String("method", r.Method))
	switch r.Method {
	case http.MethodGet:
		return s.getCategory(w, r)
	case http.MethodPut:
		return s.putCategory(w, r)
	case http.MethodDelete:
		return s.deleteCategory(w, r)
	default:
		s.logger.Error("unhandled method", zap.String("method", r.Method))
		return writeJSON(w, http.StatusInternalServerError, fmt.Errorf("unhandled method %s", r.Method), s.logger)
	}
}

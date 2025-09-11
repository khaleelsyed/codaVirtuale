package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/khaleelsyed/codaVirtuale/internal/types"
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
		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	desk, err := s.storage.GetDesk(deskID)
	if err != nil {
		if err == types.ErrnotFound {
			return writeJSON(w, http.StatusNotFound, err, s.logger)
		}
		return writeJSON(w, http.StatusBadRequest, badValidationString("desk"), s.logger)
	}

	return writeJSON(w, http.StatusOK, desk, s.logger)
}

func (s *APIServer) putDesk(w http.ResponseWriter, r *http.Request) error {
	var err error

	type DeskUpdate struct {
		CategoryID int    `json:"category_id"`
		Label      string `json:"label"`
	}

	var requestBody DeskUpdate

	idStr := mux.Vars(r)["id"]
	deskID, err := strconv.Atoi(idStr)
	if err != nil {
		errBody := "bad ID"

		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	if err = json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return writeJSON(w, http.StatusBadRequest, errBadRequestBody, s.logger)
	}

	validate := func(rB DeskUpdate) (DeskUpdate, error) {
		if rB.Label == "" && rB.CategoryID == 0 {
			s.logger.Debugw("deskUpdate check", "label", rB.Label, "CategoryID", rB.CategoryID, "id", deskID)
			return DeskUpdate{}, apiError{"Request must contain either 'category_id' or a 'label'"}
		}

		currentRow, err := s.storage.GetDesk(deskID)
		if err != nil {
			return DeskUpdate{}, err
		}

		if rB.CategoryID == 0 {
			rB.CategoryID = currentRow.CategoryID
		} else if rB.Label == "" {
			rB.Label = currentRow.Label

			if _, err = s.storage.GetCategory(rB.CategoryID); err != nil {
				if err == types.ErrnotFound {
					return DeskUpdate{}, errors.New("category not found")
				}
				return DeskUpdate{}, err
			}
		}
		return rB, nil
	}

	requestBody, err = validate(requestBody)
	if err != nil {
		if err == types.ErrnotFound || err.Error() == "category not found" {
			return writeJSON(w, http.StatusNotFound, err, s.logger)
		}

		s.logger.Tracew("failed to validate desk", "id", deskID, "error", err)
		return writeJSON(w, http.StatusBadRequest, err, s.logger)
	}

	desk, err := s.storage.UpdateDesk(deskID, struct {
		CategoryID int
		Label      string
	}(requestBody))
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, errors.New("failed to update desk"), s.logger)
	}

	return writeJSON(w, http.StatusOK, desk, s.logger)
}

func (s *APIServer) deleteDesk(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	deskID, err := strconv.Atoi(idStr)
	if err != nil {
		errBody := "bad ID"

		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	if err := s.storage.DeleteDesk(deskID); err != nil {
		errBody := badValidationString("desk")

		return writeJSON(w, http.StatusBadRequest, errBody, s.logger)
	}

	return writeJSON(w, http.StatusNoContent, nil, s.logger)
}

func (s *APIServer) createDesk(w http.ResponseWriter, r *http.Request) error {
	var err error

	type DeskCreate struct {
		Label      string `json:"label"`
		CategoryID int    `json:"category_id"`
	}

	var requestBody DeskCreate

	if err = json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return writeJSON(w, http.StatusBadRequest, errBadRequestBody, s.logger)
	}

	validate := func(rB DeskCreate) []error {
		var errs []error
		if rB.Label == "" {
			errs = append(errs, apiError{"label is required"})
		}
		if rB.CategoryID == 0 {
			errs = append(errs, apiError{"category_id is required"})
		}
		return errs
	}

	if errs := validate(requestBody); len(errs) > 0 {
		return writeJSON(w, http.StatusBadRequest, errs, s.logger)
	}

	desk, err := s.storage.CreateDesk(requestBody.Label, requestBody.CategoryID)
	if err != nil {
		if strings.Contains(err.Error(), pqForeignKeyConstraintViolation) {
			return writeJSON(w, http.StatusInternalServerError, errors.New("category_id does not exist"), s.logger)
		}
		errBody := "error creating desk"
		return writeJSON(w, http.StatusInternalServerError, errors.New(errBody), s.logger)
	}

	return writeJSON(w, http.StatusCreated, desk, s.logger)
}

func (s *APIServer) handleDesk(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.getDesk(w, r)
	case http.MethodPut:
		return s.putDesk(w, r)
	case http.MethodDelete:
		return s.deleteDesk(w, r)
	default:
		s.logger.Errorw("unhandled method", "method", r.Method)
		return writeJSON(w, http.StatusInternalServerError, fmt.Errorf("unhandled method %s", r.Method), s.logger)
	}
}

package step_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/step_actions"
	"go-upcycle_connect-backend/app/models/step_models"
	"go-upcycle_connect-backend/utils/jwt"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
	"strconv"
)

func parseStepID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, response.ErrStepNotFound, http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func findStep(w http.ResponseWriter, id int) bool {
	var s step_models.Step
	if err := s.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrStepNotFound, http.StatusNotFound)
		return false
	}
	return true
}

func GetStepsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	var s step_models.Step
	var steps []step_models.Step
	if err := s.All([]string{"id", "title", "description", "`order`", "created_at", "updated_at"}, &steps); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, steps)
}

func GetStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseStepID(w, r)
	if !ok {
		return
	}
	var s step_models.Step
	if err := s.Get([]string{"id", "title", "description", "`order`", "created_at", "updated_at"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrStepNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, s)
}

func CreateStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	var dto step_actions.CreateStepDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	s, errs := step_actions.CreateStep(dto)
	if len(errs) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, errs)
		return
	}
	if s == nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, s)
}

func UpdateStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseStepID(w, r)
	if !ok {
		return
	}
	if !findStep(w, id) {
		return
	}
	var dto step_actions.UpdateStepDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	updated, errs := step_actions.UpdateStep(id, dto)
	if len(errs) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, errs)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrStepNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, updated)
}

func DeleteStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseStepID(w, r)
	if !ok {
		return
	}
	if !findStep(w, id) {
		return
	}
	step_models.DeleteStep(id)
	response.NewSuccessMessage(w, response.SuccessDeleted)
}

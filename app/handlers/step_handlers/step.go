package step_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/step_actions"
	"go-upcycle_connect-backend/app/models/step_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func findStep(w http.ResponseWriter, id string) bool {
	var s step_models.Step
	if err := s.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrStepNotFound, http.StatusNotFound)
		return false
	}
	return true
}

func IndexStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var s step_models.Step
	var steps []step_models.Step
	columns := []string{"id", "name", "description", "image_path", "user_id", "project_id", "scheduled_at", "created_at", "updated_at"}
	if err := s.All(columns, &steps); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, steps)
}

func ShowStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrStepNotFound, http.StatusNotFound)
		return
	}
	var s step_models.Step
	columns := []string{"id", "name", "description", "image_path", "user_id", "project_id", "scheduled_at", "created_at", "updated_at"}
	if err := s.Get(columns, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrStepNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, s)
}

func StoreStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto step_models.CreateStepDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, step := step_actions.CreateStep(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if step == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]string{"id": step.Id})
}

func UpdateStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrStepNotFound, http.StatusNotFound)
		return
	}
	if !findStep(w, id) {
		return
	}
	var dto step_models.UpdateStepDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, updated := step_actions.UpdateStep(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrStepNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]string{"id": updated.Id})
}

func DeleteStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrStepNotFound, http.StatusNotFound)
		return
	}
	if !findStep(w, id) {
		return
	}
	step_models.DeleteStep(id)
	response.NewSuccessMessage(w, "Step deleted")
}

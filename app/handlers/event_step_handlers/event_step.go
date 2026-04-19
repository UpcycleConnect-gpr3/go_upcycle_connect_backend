package event_step_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/event_step_actions"
	"go-upcycle_connect-backend/app/models/event_step_models"
	"go-upcycle_connect-backend/utils/handler"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func findEventStep(w http.ResponseWriter, id int) bool {
	var s event_step_models.EventStep
	if err := s.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrEventStepNotFound, http.StatusNotFound)
		return false
	}
	return true
}

func GetEventStepsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var s event_step_models.EventStep
	var steps []event_step_models.EventStep
	if err := s.All([]string{"id", "event_id", "title", "`order`", "created_at", "updated_at"}, &steps); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, steps)
}

func GetEventStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id, ok := handler.ParsePathInt(w, r, "id", response.ErrEventStepNotFound)
	if !ok {
		return
	}
	var s event_step_models.EventStep
	if err := s.Get([]string{"id", "event_id", "title", "`order`", "created_at", "updated_at"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrEventStepNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, s)
}

func CreateEventStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto event_step_actions.CreateEventStepDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	s, errs := event_step_actions.CreateEventStep(dto)
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

func UpdateEventStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id, ok := handler.ParsePathInt(w, r, "id", response.ErrEventStepNotFound)
	if !ok {
		return
	}
	if !findEventStep(w, id) {
		return
	}
	var dto event_step_actions.UpdateEventStepDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	updated, errs := event_step_actions.UpdateEventStep(id, dto)
	if len(errs) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, errs)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrEventStepNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, updated)
}

func DeleteEventStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id, ok := handler.ParsePathInt(w, r, "id", response.ErrEventStepNotFound)
	if !ok {
		return
	}
	if !findEventStep(w, id) {
		return
	}
	event_step_models.DeleteEventStep(id)
	response.NewSuccessMessage(w, response.SuccessDeleted)
}

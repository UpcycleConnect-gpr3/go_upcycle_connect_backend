package event_step_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/event_step_actions"
	"go-upcycle_connect-backend/app/models/event_step_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
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

func IndexEventStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var s event_step_models.EventStep
	var steps []event_step_models.EventStep
	columns := []string{"id", "event_id", "title", "`order`", "created_at", "updated_at"}
	if err := s.All(columns, &steps); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, steps)
}

func ShowEventStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var s event_step_models.EventStep
	columns := []string{"id", "event_id", "title", "`order`", "created_at", "updated_at"}
	if err := s.Get(columns, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrEventStepNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, s)
}

func StoreEventStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto event_step_models.CreateEventStepDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, eventStep := event_step_actions.CreateEventStep(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if eventStep == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": eventStep.Id})
}

func UpdateEventStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findEventStep(w, id) {
		return
	}
	var dto event_step_models.UpdateEventStepDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, updated := event_step_actions.UpdateEventStep(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrEventStepNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": updated.Id})
}

func DeleteEventStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findEventStep(w, id) {
		return
	}
	event_step_models.DeleteEventStep(id)
	response.NewSuccessMessage(w, "Event step deleted")
}

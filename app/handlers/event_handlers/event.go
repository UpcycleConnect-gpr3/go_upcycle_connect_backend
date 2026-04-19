package event_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/event_actions"
	"go-upcycle_connect-backend/app/models/event_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func findEvent(w http.ResponseWriter, id int) bool {
	var e event_models.Event
	if err := e.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusNotFound)
		return false
	}
	return true
}

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var e event_models.Event
	var events []event_models.Event
	columns := []string{"id", "title", "date", "location", "created_at", "updated_at"}
	if err := e.All(columns, &events); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, events)
}

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var e event_models.Event
	columns := []string{"id", "title", "date", "location", "created_at", "updated_at"}
	if err := e.Get(columns, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, e)
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto event_actions.CreateEventDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, event := event_actions.CreateEvent(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if event == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": event.Id})
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findEvent(w, id) {
		return
	}
	var dto event_actions.UpdateEventDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, updated := event_actions.UpdateEvent(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": updated.Id})
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findEvent(w, id) {
		return
	}
	event_models.DeleteEvent(id)
	response.NewSuccessMessage(w, "Event deleted")
}

func GetEventStepsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findEvent(w, id) {
		return
	}
	steps := event_models.GetEventSteps(id)
	response.NewSuccessData(w, steps)
}

func CreateEventStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findEvent(w, id) {
		return
	}
	var dto event_actions.CreateEventStepDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, eventStep := event_actions.CreateEventStep(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if eventStep == nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": eventStep.Id})
}

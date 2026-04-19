package event_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/event_actions"
	"go-upcycle_connect-backend/app/models/event_models"
	"go-upcycle_connect-backend/utils/jwt"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
	"strconv"
)

func parseEventID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

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
	if !jwt.Auth(w, r) {
		return
	}
	var e event_models.Event
	var events []event_models.Event
	if err := e.All([]string{"id", "title", "date", "location", "created_at", "updated_at"}, &events); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, events)
}

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseEventID(w, r)
	if !ok {
		return
	}
	var e event_models.Event
	if err := e.Get([]string{"id", "title", "date", "location", "created_at", "updated_at"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, e)
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	var dto event_actions.CreateEventDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	e, errs := event_actions.CreateEvent(dto)
	if len(errs) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, errs)
		return
	}
	if e == nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, e)
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseEventID(w, r)
	if !ok {
		return
	}
	if !findEvent(w, id) {
		return
	}
	var dto event_actions.UpdateEventDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	updated, errs := event_actions.UpdateEvent(id, dto)
	if len(errs) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, errs)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, updated)
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseEventID(w, r)
	if !ok {
		return
	}
	if !findEvent(w, id) {
		return
	}
	event_models.DeleteEvent(id)
	response.NewSuccessMessage(w, response.SuccessDeleted)
}

func GetEventStepsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseEventID(w, r)
	if !ok {
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
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseEventID(w, r)
	if !ok {
		return
	}
	if !findEvent(w, id) {
		return
	}
	var dto event_actions.CreateEventStepDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	s, errs := event_actions.CreateEventStep(id, dto)
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

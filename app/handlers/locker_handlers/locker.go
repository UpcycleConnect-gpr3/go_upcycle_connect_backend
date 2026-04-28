package locker_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/locker_actions"
	"go-upcycle_connect-backend/app/models/locker_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func IndexLockerHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var locker locker_models.Locker
	var lockers []locker_models.Locker
	columns := []string{"id", "name", "street", "city", "zip_code", "created_at", "updated_at"}
	if err := locker.All(columns, &lockers); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, lockers)
}

func ShowLockerHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrLockerNotFound, http.StatusNotFound)
		return
	}
	var locker locker_models.Locker
	columns := []string{"id", "name", "street", "city", "zip_code", "created_at", "updated_at"}
	if err := locker.Get(columns, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrLockerNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, locker)
}

func StoreLockerHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto locker_models.CreateLockerDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, locker := locker_actions.CreateLocker(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if locker == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]string{"id": locker.Id})
}

func UpdateLockerHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrLockerNotFound, http.StatusNotFound)
		return
	}
	var locker locker_models.Locker
	if err := locker.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrLockerNotFound, http.StatusNotFound)
		return
	}
	var dto locker_models.UpdateLockerDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, updated := locker_actions.UpdateLocker(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrLockerNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]string{"id": updated.Id})
}

func DeleteLockerHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrLockerNotFound, http.StatusNotFound)
		return
	}
	var locker locker_models.Locker
	if err := locker.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrLockerNotFound, http.StatusNotFound)
		return
	}
	locker_models.DeleteLocker(id)
	response.NewSuccessMessage(w, "Locker deleted successfully")
}

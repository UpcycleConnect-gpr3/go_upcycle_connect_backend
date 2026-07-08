package appointment_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/appointment_actions"
	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/appointment_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
	"strconv"
)

func findOwnedAppointment(w http.ResponseWriter, r *http.Request) *appointment_models.Appointment {
	id, err := strconv.Atoi(request.Request(r, "id").Value())
	if err != nil {
		response.NewErrorMessage(w, response.ErrAppointmentNotFound, http.StatusNotFound)
		return nil
	}

	var appointment appointment_models.Appointment
	if err := appointment.Get([]string{"id", "user_id"}, "id = ?", id); err != nil {
		response.NewErrorMessage(w, response.ErrAppointmentNotFound, http.StatusNotFound)
		return nil
	}

	if appointment.UserId != auth_middleware.GetUserId(r.Context()) {
		response.NewErrorMessage(w, response.ErrAppointmentNotFound, http.StatusNotFound)
		return nil
	}

	return &appointment
}

func IndexAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())
	appointments := appointment_models.GetUserAppointments(userId)
	response.NewSuccessData(w, appointments)
}

func StoreAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())

	var dto appointment_models.CreateAppointmentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, appointment := appointment_actions.CreateAppointment(userId, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if appointment == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}

	response.NewSuccessData(w, appointment)
}

func UpdateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	appointment := findOwnedAppointment(w, r)
	if appointment == nil {
		return
	}

	var dto appointment_models.UpdateAppointmentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, updated := appointment_actions.UpdateAppointment(appointment.Id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}

	response.NewSuccessData(w, map[string]int{"id": updated.Id})
}

func DeleteAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	appointment := findOwnedAppointment(w, r)
	if appointment == nil {
		return
	}

	appointment_models.DeleteAppointment(appointment.Id)
	response.NewSuccessMessage(w, "Appointment deleted successfully")
}

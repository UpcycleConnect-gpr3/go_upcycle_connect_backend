package appointment_actions

import (
	"go-upcycle_connect-backend/app/models/appointment_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateAppointment(userId string, dto appointment_models.CreateAppointmentDTO) ([]rules.ValidationError, *appointment_models.Appointment) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMinLength(dto.StartsAt, 1, "starts_at", &errs)
	rules.StringMinLength(dto.EndsAt, 1, "ends_at", &errs)

	if dto.Kind == "" {
		dto.Kind = "event"
	}

	if len(errs) > 0 {
		return errs, nil
	}

	appointment := appointment_models.CreateAppointment(userId, dto)

	return nil, appointment
}

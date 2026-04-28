package event_step_actions

import (
	"go-upcycle_connect-backend/app/models/event_step_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateEventStep(dto event_step_models.CreateEventStepDTO) ([]rules.ValidationError, *event_step_models.EventStep) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	eventStep := event_step_models.CreateEventStep(dto)

	return nil, eventStep
}

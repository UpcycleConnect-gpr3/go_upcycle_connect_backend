package event_step_actions

import (
	"go-upcycle_connect-backend/app/models/event_step_models"
	"go-upcycle_connect-backend/utils/rules"
)

func UpdateEventStep(id int, dto event_step_models.UpdateEventStepDTO) ([]rules.ValidationError, *event_step_models.EventStep) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	eventStep := event_step_models.UpdateEventStep(id, dto)

	return nil, eventStep
}

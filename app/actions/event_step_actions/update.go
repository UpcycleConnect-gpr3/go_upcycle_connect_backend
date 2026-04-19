package event_step_actions

import (
	"go-upcycle_connect-backend/app/models/event_step_models"
	"go-upcycle_connect-backend/utils/rules"
)

type UpdateEventStepDTO struct {
	Title string `json:"title"`
	Order int    `json:"order"`
}

func UpdateEventStep(id int, dto UpdateEventStepDTO) ([]rules.ValidationError, *event_step_models.EventStep) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	eventStep := event_step_models.UpdateEventStep(id, event_step_models.UpdateEventStepDTO{
		Title: dto.Title,
		Order: dto.Order,
	})

	return nil, eventStep
}

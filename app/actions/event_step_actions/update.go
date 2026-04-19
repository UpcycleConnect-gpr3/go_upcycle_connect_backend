package event_step_actions

import (
	"go-upcycle_connect-backend/app/models/event_step_models"
	"go-upcycle_connect-backend/utils/rules"
)

type UpdateEventStepDTO struct {
	Title string `json:"title"`
	Order int    `json:"order"`
}

func validateUpdate(dto UpdateEventStepDTO) []rules.ValidationError {
	var errs []rules.ValidationError
	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)
	return errs
}

func UpdateEventStep(id int, dto UpdateEventStepDTO) (*event_step_models.EventStep, []rules.ValidationError) {
	errs := validateUpdate(dto)
	if len(errs) > 0 {
		return nil, errs
	}
	s := event_step_models.UpdateEventStep(id, event_step_models.UpdateEventStepDTO{
		Title: dto.Title,
		Order: dto.Order,
	})
	return s, nil
}

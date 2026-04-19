package event_actions

import (
	"go-upcycle_connect-backend/app/models/event_models"
	"go-upcycle_connect-backend/utils/rules"
)

type CreateEventStepDTO struct {
	Title string `json:"title"`
	Order int    `json:"order"`
}

func CreateEventStep(eventID int, dto CreateEventStepDTO) ([]rules.ValidationError, *event_models.EventStepSummary) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	eventStepSummary := event_models.CreateEventStep(eventID, dto.Title, dto.Order)
	return nil, eventStepSummary
}

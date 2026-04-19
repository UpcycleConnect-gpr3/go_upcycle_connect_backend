package event_actions

import (
	"go-upcycle_connect-backend/app/models/event_models"
	"go-upcycle_connect-backend/utils/rules"
)

type CreateEventStepDTO struct {
	Title string `json:"title"`
	Order int    `json:"order"`
}

func validateCreateStep(dto CreateEventStepDTO) []rules.ValidationError {
	var errs []rules.ValidationError
	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)
	return errs
}

func CreateEventStep(eventID int, dto CreateEventStepDTO) (*event_models.EventStepSummary, []rules.ValidationError) {
	errs := validateCreateStep(dto)
	if len(errs) > 0 {
		return nil, errs
	}
	s := event_models.CreateEventStep(eventID, dto.Title, dto.Order)
	return s, nil
}

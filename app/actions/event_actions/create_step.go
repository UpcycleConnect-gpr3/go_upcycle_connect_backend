package event_actions

import (
	"go-upcycle_connect-backend/app/models/event_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateEventStep(eventID int, name, description, imagePath string, scheduledAt string) ([]rules.ValidationError, *event_models.EventStepSummary) {
	var errs []rules.ValidationError

	rules.StringMinLength(name, 1, "name", &errs)
	rules.StringMaxLength(name, 255, "name", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	eventStepSummary := event_models.CreateEventStep(eventID, name, description, imagePath, scheduledAt)
	return nil, eventStepSummary
}

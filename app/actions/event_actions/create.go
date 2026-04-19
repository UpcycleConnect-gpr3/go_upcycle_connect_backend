package event_actions

import (
	"go-upcycle_connect-backend/app/models/event_models"
	"go-upcycle_connect-backend/utils/rules"
)

type CreateEventDTO struct {
	Title    string `json:"title"`
	Date     string `json:"date"`
	Location string `json:"location"`
}

func CreateEvent(dto CreateEventDTO) ([]rules.ValidationError, *event_models.Event) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	event := event_models.CreateEvent(event_models.CreateEventDTO{
		Title:    dto.Title,
		Date:     dto.Date,
		Location: dto.Location,
	})

	return nil, event
}

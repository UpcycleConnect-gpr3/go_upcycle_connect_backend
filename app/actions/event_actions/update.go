package event_actions

import (
	"go-upcycle_connect-backend/app/models/event_models"
	"go-upcycle_connect-backend/utils/rules"
)

func UpdateEvent(id int, dto event_models.UpdateEventDTO) ([]rules.ValidationError, *event_models.Event) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	event := event_models.UpdateEvent(id, dto)

	return nil, event
}

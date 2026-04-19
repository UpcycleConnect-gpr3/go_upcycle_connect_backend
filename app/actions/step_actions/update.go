package step_actions

import (
	"go-upcycle_connect-backend/app/models/step_models"
	"go-upcycle_connect-backend/utils/rules"
)

type UpdateStepDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

func validateUpdate(dto UpdateStepDTO) []rules.ValidationError {
	var errs []rules.ValidationError
	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)
	return errs
}

func UpdateStep(id int, dto UpdateStepDTO) (*step_models.Step, []rules.ValidationError) {
	errs := validateUpdate(dto)
	if len(errs) > 0 {
		return nil, errs
	}
	s := step_models.UpdateStep(id, step_models.UpdateStepDTO{
		Title:       dto.Title,
		Description: dto.Description,
		Order:       dto.Order,
	})
	return s, nil
}

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

func UpdateStep(id int, dto UpdateStepDTO) ([]rules.ValidationError, *step_models.Step) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	step := step_models.UpdateStep(id, step_models.UpdateStepDTO{
		Title:       dto.Title,
		Description: dto.Description,
		Order:       dto.Order,
	})

	return nil, step
}

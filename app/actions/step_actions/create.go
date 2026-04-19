package step_actions

import (
	"go-upcycle_connect-backend/app/models/step_models"
	"go-upcycle_connect-backend/utils/rules"
)

type CreateStepDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

func validateCreate(dto CreateStepDTO) []rules.ValidationError {
	var errs []rules.ValidationError
	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)
	return errs
}

func CreateStep(dto CreateStepDTO) (*step_models.Step, []rules.ValidationError) {
	errs := validateCreate(dto)
	if len(errs) > 0 {
		return nil, errs
	}
	s := step_models.CreateStep(step_models.CreateStepDTO{
		Title:       dto.Title,
		Description: dto.Description,
		Order:       dto.Order,
	})
	return s, nil
}

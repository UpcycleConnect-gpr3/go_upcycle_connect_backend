package project_actions

import (
	"go-upcycle_connect-backend/app/models/project_models"
	"go-upcycle_connect-backend/utils/rules"
)

type CreateProjectStepDTO struct {
	Title string `json:"title"`
	Order int    `json:"order"`
}

func validateCreateStep(dto CreateProjectStepDTO) []rules.ValidationError {
	var errs []rules.ValidationError
	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)
	return errs
}

func CreateProjectStep(projectID int, dto CreateProjectStepDTO) (*project_models.StepSummary, []rules.ValidationError) {
	errs := validateCreateStep(dto)
	if len(errs) > 0 {
		return nil, errs
	}
	s := project_models.CreateProjectStep(projectID, dto.Title, dto.Order)
	return s, nil
}

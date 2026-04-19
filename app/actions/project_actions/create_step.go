package project_actions

import (
	"go-upcycle_connect-backend/app/models/project_models"
	"go-upcycle_connect-backend/utils/rules"
)

type CreateProjectStepDTO struct {
	Title string `json:"title"`
	Order int    `json:"order"`
}

func CreateProjectStep(projectID int, dto CreateProjectStepDTO) ([]rules.ValidationError, *project_models.StepSummary) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	stepSummary := project_models.CreateProjectStep(projectID, dto.Title, dto.Order)
	return nil, stepSummary
}

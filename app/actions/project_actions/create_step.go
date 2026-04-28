package project_actions

import (
	"go-upcycle_connect-backend/app/models/project_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateProjectStep(projectID int, name, description, imagePath string, scheduledAt string) ([]rules.ValidationError, *project_models.StepSummary) {
	var errs []rules.ValidationError

	rules.StringMinLength(name, 1, "name", &errs)
	rules.StringMaxLength(name, 255, "name", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	stepSummary := project_models.CreateProjectStep(projectID, name, description, imagePath, scheduledAt)
	return nil, stepSummary
}

package project_actions

import (
	"go-upcycle_connect-backend/app/models/project_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateProject(dto project_models.CreateProjectDTO) ([]rules.ValidationError, *project_models.Project) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	project := project_models.CreateProject(dto)

	return nil, project
}

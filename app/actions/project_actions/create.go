package project_actions

import (
	"go-upcycle_connect-backend/app/models/project_models"
	"go-upcycle_connect-backend/utils/rules"
)

type CreateProjectDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateProject(dto CreateProjectDTO) ([]rules.ValidationError, *project_models.Project) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	project := project_models.CreateProject(project_models.CreateProjectDTO{
		Name:        dto.Name,
		Description: dto.Description,
	})

	return nil, project
}

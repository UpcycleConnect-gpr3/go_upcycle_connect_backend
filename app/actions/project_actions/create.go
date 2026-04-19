package project_actions

import (
	"go-upcycle_connect-backend/app/models/project_models"
	"go-upcycle_connect-backend/utils/rules"
)

type CreateProjectDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func validateCreate(dto CreateProjectDTO) []rules.ValidationError {
	var errs []rules.ValidationError
	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)
	return errs
}

func CreateProject(dto CreateProjectDTO) (*project_models.Project, []rules.ValidationError) {
	errs := validateCreate(dto)
	if len(errs) > 0 {
		return nil, errs
	}
	p := project_models.CreateProject(project_models.CreateProjectDTO{
		Name:        dto.Name,
		Description: dto.Description,
	})
	return p, nil
}

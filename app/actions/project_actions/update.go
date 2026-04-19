package project_actions

import (
	"go-upcycle_connect-backend/app/models/project_models"
	"go-upcycle_connect-backend/utils/rules"
)

type UpdateProjectDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func validateUpdate(dto UpdateProjectDTO) []rules.ValidationError {
	var errs []rules.ValidationError
	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)
	return errs
}

func UpdateProject(id int, dto UpdateProjectDTO) (*project_models.Project, []rules.ValidationError) {
	errs := validateUpdate(dto)
	if len(errs) > 0 {
		return nil, errs
	}
	p := project_models.UpdateProject(id, project_models.UpdateProjectDTO{
		Name:        dto.Name,
		Description: dto.Description,
	})
	return p, nil
}

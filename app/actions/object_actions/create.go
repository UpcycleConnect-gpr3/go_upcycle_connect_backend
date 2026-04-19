package object_actions

import (
	"go-upcycle_connect-backend/app/models/object_models"
	"go-upcycle_connect-backend/utils/rules"
)

type CreateObjectDTO struct {
	Name        string `json:"name"`
	Material    string `json:"material"`
	Condition   string `json:"condition"`
	Description string `json:"description"`
}

func validateCreate(dto CreateObjectDTO) []rules.ValidationError {
	var errs []rules.ValidationError
	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)
	return errs
}

func CreateObject(dto CreateObjectDTO) (*object_models.Object, []rules.ValidationError) {
	errs := validateCreate(dto)
	if len(errs) > 0 {
		return nil, errs
	}
	obj := object_models.CreateObject(object_models.CreateObjectDTO{
		Name:        dto.Name,
		Material:    dto.Material,
		Condition:   dto.Condition,
		Description: dto.Description,
	})
	return obj, nil
}

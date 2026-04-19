package object_actions

import (
	"go-upcycle_connect-backend/app/models/object_models"
	"go-upcycle_connect-backend/utils/rules"
)

type UpdateObjectDTO struct {
	Name        string `json:"name"`
	Material    string `json:"material"`
	Condition   string `json:"condition"`
	Description string `json:"description"`
}

func validateUpdate(dto UpdateObjectDTO) []rules.ValidationError {
	var errs []rules.ValidationError
	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)
	return errs
}

func UpdateObject(id int, dto UpdateObjectDTO) (*object_models.Object, []rules.ValidationError) {
	errs := validateUpdate(dto)
	if len(errs) > 0 {
		return nil, errs
	}
	obj := object_models.UpdateObject(id, object_models.UpdateObjectDTO{
		Name:        dto.Name,
		Material:    dto.Material,
		Condition:   dto.Condition,
		Description: dto.Description,
	})
	return obj, nil
}

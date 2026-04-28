package object_actions

import (
	"go-upcycle_connect-backend/app/models/object_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateObject(dto object_models.CreateObjectDTO) ([]rules.ValidationError, *object_models.Object) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	obj := object_models.CreateObject(dto)

	return nil, obj
}

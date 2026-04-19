package delivery_method_actions

import (
	"go-upcycle_connect-backend/app/models/delivery_method_models"
	"go-upcycle_connect-backend/utils/rules"
)

type UpdateDeliveryMethodDTO struct {
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

func validateUpdate(dto UpdateDeliveryMethodDTO) []rules.ValidationError {
	var errs []rules.ValidationError
	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)
	return errs
}

func UpdateDeliveryMethod(id int, dto UpdateDeliveryMethodDTO) (*delivery_method_models.DeliveryMethod, []rules.ValidationError) {
	errs := validateUpdate(dto)
	if len(errs) > 0 {
		return nil, errs
	}
	dm := delivery_method_models.UpdateDeliveryMethod(id, delivery_method_models.UpdateDeliveryMethodDTO{
		Name: dto.Name,
		Cost: dto.Cost,
	})
	return dm, nil
}

package delivery_method_actions

import (
	"go-upcycle_connect-backend/app/models/delivery_method_models"
	"go-upcycle_connect-backend/utils/rules"
)

type CreateDeliveryMethodDTO struct {
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

func validateCreate(dto CreateDeliveryMethodDTO) []rules.ValidationError {
	var errs []rules.ValidationError
	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)
	return errs
}

func CreateDeliveryMethod(dto CreateDeliveryMethodDTO) (*delivery_method_models.DeliveryMethod, []rules.ValidationError) {
	errs := validateCreate(dto)
	if len(errs) > 0 {
		return nil, errs
	}
	dm := delivery_method_models.CreateDeliveryMethod(delivery_method_models.CreateDeliveryMethodDTO{
		Name: dto.Name,
		Cost: dto.Cost,
	})
	return dm, nil
}

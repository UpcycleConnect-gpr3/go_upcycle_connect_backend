package delivery_method_actions

import (
	"go-upcycle_connect-backend/app/models/delivery_method_models"
	"go-upcycle_connect-backend/utils/rules"
)

func UpdateDeliveryMethod(id int, dto delivery_method_models.UpdateDeliveryMethodDTO) ([]rules.ValidationError, *delivery_method_models.DeliveryMethod) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	deliveryMethod := delivery_method_models.UpdateDeliveryMethod(id, dto)

	return nil, deliveryMethod
}

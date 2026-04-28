package order_delivery_method_actions

import (
	"go-upcycle_connect-backend/app/models/order_delivery_method_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateOrderDeliveryMethod(dto order_delivery_method_models.CreateOrderDeliveryMethodDTO) ([]rules.ValidationError, *order_delivery_method_models.OrderDeliveryMethod) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.OrderId, 1, "order_id", &errs)
	rules.IntMinLength(dto.DeliveryMethodId, 1, "delivery_method_id", &errs)
	if dto.Price < 0 {
		errs = append(errs, rules.ValidationError{Field: "price", Message: "price must be at least 0"})
	}

	if len(errs) > 0 {
		return errs, nil
	}

	odm := order_delivery_method_models.CreateOrderDeliveryMethod(dto)

	return nil, odm
}

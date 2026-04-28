package order_delivery_method_actions

import (
	"go-upcycle_connect-backend/app/models/order_delivery_method_models"
	"go-upcycle_connect-backend/utils/rules"
)

func UpdateOrderDeliveryMethod(orderId string, deliveryMethodId int, dto order_delivery_method_models.UpdateOrderDeliveryMethodDTO) ([]rules.ValidationError, *order_delivery_method_models.OrderDeliveryMethod) {
	var errs []rules.ValidationError

	if dto.Price < 0 {
		errs = append(errs, rules.ValidationError{Field: "price", Message: "price must be at least 0"})
	}

	if len(errs) > 0 {
		return errs, nil
	}

	odm := order_delivery_method_models.UpdateOrderDeliveryMethod(orderId, deliveryMethodId, dto)

	return nil, odm
}

package order_actions

import (
	"go-upcycle_connect-backend/app/models/order_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateOrder(dto order_models.CreateOrderDTO) ([]rules.ValidationError, *order_models.Order) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Street, 1, "street", &errs)
	rules.StringMinLength(dto.City, 1, "city", &errs)
	rules.StringMinLength(dto.ZipCode, 1, "zip_code", &errs)
	rules.StringMinLength(dto.UserId, 1, "user_id", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	order := order_models.CreateOrder(dto)

	return nil, order
}

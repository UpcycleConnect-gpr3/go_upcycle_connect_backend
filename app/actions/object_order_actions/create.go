package object_order_actions

import (
	"go-upcycle_connect-backend/app/models/object_order_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateObjectOrder(dto object_order_models.CreateObjectOrderDTO) ([]rules.ValidationError, *object_order_models.ObjectOrder) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.ObjectId, 1, "object_id", &errs)
	rules.StringMinLength(dto.OrderId, 1, "order_id", &errs)
	rules.IntMinLength(dto.Amount, 1, "amount", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	oo := object_order_models.CreateObjectOrder(dto)

	return nil, oo
}

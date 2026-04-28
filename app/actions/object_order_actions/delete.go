package object_order_actions

import (
	"go-upcycle_connect-backend/app/models/object_order_models"
	"go-upcycle_connect-backend/utils/rules"
)

func DeleteObjectOrder(id int) ([]rules.ValidationError, bool) {
	var errs []rules.ValidationError

	rules.IntMinLength(id, 1, "id", &errs)

	if len(errs) > 0 {
		return errs, false
	}

	var oo object_order_models.ObjectOrder
	if err := oo.Get([]string{"id"}, "id", id); err != nil {
		return []rules.ValidationError{{Field: "id", Message: "ObjectOrder not found"}}, false
	}

	object_order_models.DeleteObjectOrder(id)
	return nil, true
}

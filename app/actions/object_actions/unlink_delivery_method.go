package object_actions

import (
	"go-upcycle_connect-backend/app/models/object_models"
)

func UnlinkDeliveryMethod(objectID, deliveryMethodID int) {
	object_models.UnlinkDeliveryMethod(objectID, deliveryMethodID)
}

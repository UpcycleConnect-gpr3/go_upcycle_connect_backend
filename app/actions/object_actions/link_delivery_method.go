package object_actions

import (
	"go-upcycle_connect-backend/app/models/object_models"
)

func LinkDeliveryMethod(objectID string, deliveryMethodID int) {
	object_models.LinkDeliveryMethod(objectID, deliveryMethodID)
}

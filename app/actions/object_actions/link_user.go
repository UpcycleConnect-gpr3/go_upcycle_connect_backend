package object_actions

import (
	"go-upcycle_connect-backend/app/models/object_models"
)

func LinkUser(objectID int, userID string) {
	object_models.LinkUser(objectID, userID)
}

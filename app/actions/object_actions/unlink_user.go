package object_actions

import (
	"go-upcycle_connect-backend/app/models/object_models"
)

func UnlinkUser(objectID int, userID string) {
	object_models.UnlinkUser(objectID, userID)
}

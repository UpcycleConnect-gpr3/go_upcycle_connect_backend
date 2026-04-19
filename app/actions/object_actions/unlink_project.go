package object_actions

import (
	"go-upcycle_connect-backend/app/models/object_models"
)

func UnlinkProject(objectID, projectID int) {
	object_models.UnlinkProject(objectID, projectID)
}

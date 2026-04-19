package object_actions

import (
	"go-upcycle_connect-backend/app/models/object_models"
)

func LinkProject(objectID, projectID int) {
	object_models.LinkProject(objectID, projectID)
}

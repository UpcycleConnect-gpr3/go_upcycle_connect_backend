package project_actions

import (
	"go-upcycle_connect-backend/app/models/project_models"
)

func LinkObject(projectID, objectID int) {
	project_models.LinkObject(projectID, objectID)
}

package project_actions

import (
	"go-upcycle_connect-backend/app/models/project_models"
)

func LinkObject(projectID int, objectID string, userID string) {
	project_models.LinkObject(projectID, objectID, userID)
}

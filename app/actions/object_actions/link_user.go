package object_actions

import (
	"github.com/google/uuid"
	"go-upcycle_connect-backend/app/models/object_models"
)

func LinkUser(objectID int, userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return err
	}
	return object_models.LinkUser(objectID, userID)
}

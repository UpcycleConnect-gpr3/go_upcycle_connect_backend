package object_models

import (
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"
)

func SetStatusAndScore(id, status string, score float64) error {
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET status = ?, score = ?, updated_at = NOW() WHERE id = ?",
		status, score, id,
	)
	if err != nil {
		log.Database("SET OBJECT STATUS+SCORE", err)
	}
	return err
}

func SetStatusAndOwner(id, status, userId string) error {
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET status = ?, user_id = ?, updated_at = NOW() WHERE id = ?",
		status, userId, id,
	)
	if err != nil {
		log.Database("SET OBJECT STATUS+OWNER", err)
	}
	return err
}

package object_models

import (
	"go-upcycle_connect-backend/var/database"
	"go-upcycle_connect-backend/utils/log"
)

// SetStatusAndScore met l'objet dans un locker : status + score CO2 calcule.
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

// SetStatusAndOwner transfere la propriete de l'objet lors de la recuperation.
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

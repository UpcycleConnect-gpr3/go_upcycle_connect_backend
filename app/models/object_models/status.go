package object_models

import (
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"
)

func SetAdValidated(id string, validated bool) error {
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET is_ad_validated = ?, updated_at = NOW() WHERE id = ?",
		validated, id,
	)
	if err != nil {
		log.Database("SET OBJECT AD VALIDATED", err)
	}
	return err
}

func AllValidated(columns []string, dest *[]Object) error {
	cols := "*"
	if len(columns) > 0 {
		cols = ""
		for i, c := range columns {
			if i > 0 {
				cols += ", "
			}
			cols += c
		}
	}
	return database.UpcycleConnect.Select(dest,
		"SELECT "+cols+" FROM "+TABLE+" WHERE is_ad_validated = 1 ORDER BY created_at DESC")
}

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

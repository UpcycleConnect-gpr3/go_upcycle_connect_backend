package activity_models

import (
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"
)

const TABLE = "ACTIVITY_LOGS"

type ActivityLog struct {
	Id        int    `db:"id" json:"id"`
	UserId    string `db:"user_id" json:"user_id"`
	Action    string `db:"action" json:"action"`
	Entity    string `db:"entity" json:"entity"`
	EntityId  string `db:"entity_id" json:"entity_id"`
	Detail    string `db:"detail" json:"detail"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

func Record(userId, action, entity, entityId, detail string) {
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (user_id, action, entity, entity_id, detail) VALUES (?, ?, ?, ?, ?)",
		userId, action, entity, entityId, detail,
	)
	if err != nil {
		log.Database("INSERT ACTIVITY LOG", err)
	}
}

func GetAll(limit int) []ActivityLog {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	logs := []ActivityLog{}
	err := database.UpcycleConnect.Select(&logs,
		"SELECT id, user_id, action, entity, entity_id, detail, created_at "+
			"FROM "+TABLE+" ORDER BY created_at DESC LIMIT ?",
		limit,
	)
	if err != nil {
		log.Database("SELECT ACTIVITY LOGS", err)
	}
	return logs
}

package notification_models

import (
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"
)

const TABLE = "NOTIFICATIONS"

type Notification struct {
	Id        int    `db:"id" json:"id"`
	UserId    string `db:"user_id" json:"user_id"`
	Title     string `db:"title" json:"title"`
	Body      string `db:"body" json:"body"`
	IsRead    bool   `db:"is_read" json:"is_read"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

// Create envoie une notification a un utilisateur.
func Create(userId, title, body string) *Notification {
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (user_id, title, body) VALUES (?, ?, ?)",
		userId, title, body,
	)
	if err != nil {
		log.Database("INSERT NOTIFICATION", err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &Notification{Id: int(id)}
}

// GetUserNotifications renvoie les notifications de l'utilisateur (recentes d'abord).
func GetUserNotifications(userId string) []Notification {
	notifications := []Notification{}
	err := database.UpcycleConnect.Select(&notifications,
		"SELECT id, user_id, title, COALESCE(body,'') AS body, is_read, created_at, updated_at "+
			"FROM "+TABLE+" WHERE user_id = ? ORDER BY created_at DESC",
		userId,
	)
	if err != nil {
		log.Database("SELECT USER NOTIFICATIONS", err)
	}
	return notifications
}

// MarkRead marque une notification comme lue (si elle appartient a l'utilisateur).
func MarkRead(id int, userId string) error {
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET is_read = TRUE, updated_at = NOW() WHERE id = ? AND user_id = ?",
		id, userId,
	)
	if err != nil {
		log.Database("MARK NOTIFICATION READ", err)
	}
	return err
}

package ad_models

import (
	"fmt"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"
)

const TABLE = "ADS"

type Ad struct {
	Id          int    `db:"id" json:"id"`
	UserId      string `db:"user_id" json:"user_id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Budget      int    `db:"budget" json:"budget"`
	Status      string `db:"status" json:"status"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}

type CreateAdDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Budget      int    `json:"budget"`
	Status      string `json:"status"`
	UserId      string `json:"-"`
}

type UpdateAdDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Budget      int    `json:"budget"`
	Status      string `json:"status"`
}

func (ad *Ad) Get(columns []string, by string, value ...any) error {
	return db.GetQuery[Ad](database.UpcycleConnect, TABLE, ad, columns, by, value...)
}

func GetUserAds(userId string) []Ad {
	result := []Ad{}
	err := database.UpcycleConnect.Select(&result,
		"SELECT id, user_id, title, description, budget, status, created_at, updated_at FROM "+TABLE+" WHERE user_id = ? ORDER BY created_at DESC",
		userId,
	)
	if err != nil {
		log.Database("SELECT USER ADS", err)
		return []Ad{}
	}
	return result
}

func CreateAd(dto CreateAdDTO) *Ad {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Title)
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (user_id, title, description, budget, status) VALUES (?, ?, ?, ?, ?)",
		dto.UserId, dto.Title, dto.Description, dto.Budget, dto.Status,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &Ad{Id: int(id)}
}

func UpdateAd(id int, dto UpdateAdDTO) *Ad {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET title=?, description=?, budget=?, status=?, updated_at=NOW() WHERE id=?",
		dto.Title, dto.Description, dto.Budget, dto.Status, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Ad{Id: id}
}

func DeleteAd(id int) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

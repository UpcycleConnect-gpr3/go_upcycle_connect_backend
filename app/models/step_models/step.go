package step_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"time"
)

const TABLE = "STEPS"

type Step struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Order       int       `json:"order"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type CreateStepDTO struct {
	Title       string
	Description string
	Order       int
}

type UpdateStepDTO struct {
	Title       string
	Description string
	Order       int
}

func (step *Step) Get(columns []string, by string, value any) error {
	return db.GetQuery[Step](database.UpcycleConnect, TABLE, columns, by, value, step)
}

func (step *Step) All(columns []string, dest *[]Step) error {
	return db.AllQuery[Step](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateStep(dto CreateStepDTO) *Step {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Title)
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (title, description, `order`) VALUES (?, ?, ?)",
		dto.Title, dto.Description, dto.Order,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &Step{Id: int(id)}
}

func UpdateStep(id int, dto UpdateStepDTO) *Step {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET title=?, description=?, `order`=? WHERE id=?",
		dto.Title, dto.Description, dto.Order, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Step{Id: id}
}

func DeleteStep(id int) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

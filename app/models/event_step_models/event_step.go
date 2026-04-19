package event_step_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
)

const TABLE = "EVENT_STEPS"

type EventStep struct {
	Id        int    `db:"id" json:"id"`
	EventId   int    `db:"event_id" json:"event_id"`
	Title     string `db:"title" json:"title"`
	Order     int    `db:"order" json:"order"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type CreateEventStepDTO struct {
	EventId int
	Title   string
	Order   int
}

type UpdateEventStepDTO struct {
	Title string
	Order int
}

func (eventStep *EventStep) Get(columns []string, by string, value any) error {
	return db.GetQuery[EventStep](database.UpcycleConnect, TABLE, columns, by, value, eventStep)
}

func (eventStep *EventStep) All(columns []string, dest *[]EventStep) error {
	return db.AllQuery[EventStep](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateEventStep(dto CreateEventStepDTO) *EventStep {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Title)
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (event_id, title, `order`) VALUES (?, ?, ?)",
		dto.EventId, dto.Title, dto.Order,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &EventStep{Id: int(id)}
}

func UpdateEventStep(id int, dto UpdateEventStepDTO) *EventStep {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET title=?, `order`=? WHERE id=?",
		dto.Title, dto.Order, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &EventStep{Id: id}
}

func DeleteEventStep(id int) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

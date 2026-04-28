package event_step_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
)

const TABLE = "EVENT_STEPS"

type EventStep struct {
	Id          int    `db:"id" json:"id"`
	EventId     int    `db:"event_id" json:"event_id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	ImagePath   string `db:"image_path" json:"image_path"`
	ScheduledAt string `db:"scheduled_at" json:"scheduled_at"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}

type CreateEventStepDTO struct {
	EventId     int
	Name        string
	Description string
	ImagePath   string
	ScheduledAt string
}

type UpdateEventStepDTO struct {
	Name        string
	Description string
	ImagePath   string
	ScheduledAt string
}

func (eventStep *EventStep) Get(columns []string, by string, value any) error {
	return db.GetQuery[EventStep](database.UpcycleConnect, TABLE, columns, by, value, eventStep)
}

func (eventStep *EventStep) All(columns []string, dest *[]EventStep) error {
	return db.AllQuery[EventStep](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateEventStep(dto CreateEventStepDTO) *EventStep {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Name)
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (event_id, name, description, image_path, scheduled_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())",
		dto.EventId, dto.Name, dto.Description, dto.ImagePath, dto.ScheduledAt,
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
		"UPDATE "+TABLE+" SET name=?, description=?, image_path=?, scheduled_at=?, updated_at=NOW() WHERE id=?",
		dto.Name, dto.Description, dto.ImagePath, dto.ScheduledAt, id,
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

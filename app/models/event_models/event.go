package event_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
)

const TABLE = "EVENTS"

type Event struct {
	Id              int    `db:"id" json:"id"`
	Name            string `db:"name" json:"name"`
	Description     string `db:"description" json:"description"`
	ImagePath       string `db:"image_path" json:"image_path"`
	StartedAt       string `db:"started_at" json:"started_at"`
	FinishedAt      string `db:"finished_at" json:"finished_at"`
	Location        string `db:"location" json:"location"`
	DeliveryMethod  string `db:"delivery_method" json:"delivery_method"`
	CreatedByUserId string `db:"created_by_user_id" json:"created_by_user_id"`
	CreatedAt       string `db:"created_at" json:"created_at"`
	UpdatedAt       string `db:"updated_at" json:"updated_at"`
}

type CreateEventDTO struct {
	Name            string
	Description     string
	ImagePath       string
	StartedAt       string
	FinishedAt      string
	Location        string
	DeliveryMethod  string
	CreatedByUserId string
}

type UpdateEventDTO struct {
	Name            string
	Description     string
	ImagePath       string
	StartedAt       string
	FinishedAt      string
	Location        string
	DeliveryMethod  string
	CreatedByUserId string
}

type EventStepSummary struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImagePath   string `json:"image_path"`
	ScheduledAt string `json:"scheduled_at"`
}

func (event *Event) Get(columns []string, by string, value any) error {
	return db.GetQuery[Event](database.UpcycleConnect, TABLE, columns, by, value, event)
}

func (event *Event) All(columns []string, dest *[]Event) error {
	return db.AllQuery[Event](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateEvent(dto CreateEventDTO) *Event {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Name)
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (name, description, image_path, started_at, finished_at, location, delivery_method, created_by_user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())",
		dto.Name, dto.Description, dto.ImagePath, dto.StartedAt, dto.FinishedAt, dto.Location, dto.DeliveryMethod, dto.CreatedByUserId,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &Event{Id: int(id)}
}

func UpdateEvent(id int, dto UpdateEventDTO) *Event {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET name=?, description=?, image_path=?, started_at=?, finished_at=?, location=?, delivery_method=?, created_by_user_id=?, updated_at=NOW() WHERE id=?",
		dto.Name, dto.Description, dto.ImagePath, dto.StartedAt, dto.FinishedAt, dto.Location, dto.DeliveryMethod, dto.CreatedByUserId, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Event{Id: id}
}

func DeleteEvent(id int) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

func GetEventSteps(eventID int) []EventStepSummary {
	result := []EventStepSummary{}
	rows, err := database.UpcycleConnect.Query(
		"SELECT id, name, description, image_path, scheduled_at FROM EVENT_STEPS WHERE event_id=?",
		eventID,
	)
	if err != nil {
		log.Database("SELECT EVENT STEPS", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		s := EventStepSummary{}
		_ = rows.Scan(&s.Id, &s.Name, &s.Description, &s.ImagePath, &s.ScheduledAt)
		result = append(result, s)
	}
	return result
}

func CreateEventStep(eventID int, name, description, imagePath, scheduledAt string) *EventStepSummary {
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO EVENT_STEPS (event_id, name, description, image_path, scheduled_at) VALUES (?, ?, ?, ?, ?)",
		eventID, name, description, imagePath, scheduledAt,
	)
	if err != nil {
		log.Database("CREATE EVENT STEP", err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &EventStepSummary{Id: int(id), Name: name, Description: description, ImagePath: imagePath, ScheduledAt: scheduledAt}
}

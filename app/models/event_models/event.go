package event_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
)

const TABLE = "EVENTS"

type Event struct {
	Id        int    `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	Date      string `db:"date" json:"date"`
	Location  string `db:"location" json:"location"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type CreateEventDTO struct {
	Title    string
	Date     string
	Location string
}

type UpdateEventDTO struct {
	Title    string
	Date     string
	Location string
}

type EventStepSummary struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Order int    `json:"order"`
}

func (event *Event) Get(columns []string, by string, value any) error {
	return db.GetQuery[Event](database.UpcycleConnect, TABLE, columns, by, value, event)
}

func (event *Event) All(columns []string, dest *[]Event) error {
	return db.AllQuery[Event](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateEvent(dto CreateEventDTO) *Event {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Title)
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (title, date, location) VALUES (?, ?, ?)",
		dto.Title, dto.Date, dto.Location,
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
		"UPDATE "+TABLE+" SET title=?, date=?, location=? WHERE id=?",
		dto.Title, dto.Date, dto.Location, id,
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
		"SELECT id, title, `order` FROM EVENT_STEPS WHERE event_id=?",
		eventID,
	)
	if err != nil {
		log.Database("SELECT EVENT STEPS", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		s := EventStepSummary{}
		_ = rows.Scan(&s.Id, &s.Title, &s.Order)
		result = append(result, s)
	}
	return result
}

func CreateEventStep(eventID int, title string, order int) *EventStepSummary {
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO EVENT_STEPS (event_id, title, `order`) VALUES (?, ?, ?)",
		eventID, title, order,
	)
	if err != nil {
		log.Database("CREATE EVENT STEP", err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &EventStepSummary{Id: int(id), Title: title, Order: order}
}

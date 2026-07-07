package appointment_models

import (
	"fmt"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"
)

const TABLE = "APPOINTMENTS"

type Appointment struct {
	Id        int    `db:"id" json:"id"`
	UserId    string `db:"user_id" json:"user_id"`
	Title     string `db:"title" json:"title"`
	Kind      string `db:"kind" json:"kind"`
	Location  string `db:"location" json:"location"`
	StartsAt  string `db:"starts_at" json:"starts_at"`
	EndsAt    string `db:"ends_at" json:"ends_at"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type CreateAppointmentDTO struct {
	Title    string `json:"title"`
	Kind     string `json:"kind"`
	Location string `json:"location"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

type UpdateAppointmentDTO struct {
	Title    string `json:"title"`
	Kind     string `json:"kind"`
	Location string `json:"location"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

func (m *Appointment) Get(columns []string, by string, value ...any) error {
	return db.GetQuery[Appointment](database.UpcycleConnect, TABLE, m, columns, by, value...)
}

func (m *Appointment) All(columns []string, dest *[]Appointment) error {
	return db.AllQuery[Appointment](database.UpcycleConnect, TABLE, columns, dest)
}

func GetUserAppointments(userId string) []Appointment {
	action := fmt.Sprintf("SELECT %s WHERE user_id: %s", TABLE, userId)
	appointments := []Appointment{}
	err := database.UpcycleConnect.Select(&appointments,
		"SELECT id, user_id, title, kind, location, starts_at, ends_at, created_at, updated_at FROM "+TABLE+" WHERE user_id = ? ORDER BY starts_at ASC",
		userId,
	)
	if err != nil {
		log.Database(action, err)
		return []Appointment{}
	}
	return appointments
}

func CreateAppointment(userId string, dto CreateAppointmentDTO) *Appointment {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Title)
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (user_id, title, kind, location, starts_at, ends_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())",
		userId, dto.Title, dto.Kind, dto.Location, dto.StartsAt, dto.EndsAt,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &Appointment{
		Id:       int(id),
		UserId:   userId,
		Title:    dto.Title,
		Kind:     dto.Kind,
		Location: dto.Location,
		StartsAt: dto.StartsAt,
		EndsAt:   dto.EndsAt,
	}
}

func UpdateAppointment(id int, dto UpdateAppointmentDTO) *Appointment {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET title=?, kind=?, location=?, starts_at=?, ends_at=?, updated_at=NOW() WHERE id=?",
		dto.Title, dto.Kind, dto.Location, dto.StartsAt, dto.EndsAt, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Appointment{Id: id}
}

func DeleteAppointment(id int) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

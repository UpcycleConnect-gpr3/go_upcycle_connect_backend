package locker_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"

	"github.com/google/uuid"
)

const TABLE = "LOCKERS"

type Locker struct {
	Id        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Street    string `db:"street" json:"street"`
	City      string `db:"city" json:"city"`
	ZipCode   string `db:"zip_code" json:"zip_code"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type CreateLockerDTO struct {
	Name    string
	Street  string
	City    string
	ZipCode string
}

type UpdateLockerDTO struct {
	Name    string
	Street  string
	City    string
	ZipCode string
}

func (m *Locker) Get(columns []string, by string, value any) error {
	return db.GetQuery[Locker](database.UpcycleConnect, TABLE, columns, by, value, m)
}

func (m *Locker) All(columns []string, dest *[]Locker) error {
	return db.AllQuery[Locker](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateLocker(dto CreateLockerDTO) *Locker {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Name)
	lockerId := uuid.New().String()
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (id, name, street, city, zip_code, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())",
		lockerId, dto.Name, dto.Street, dto.City, dto.ZipCode,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Locker{Id: lockerId}
}

func UpdateLocker(id string, dto UpdateLockerDTO) *Locker {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %s", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET name=?, street=?, city=?, zip_code=?, updated_at=NOW() WHERE id=?",
		dto.Name, dto.Street, dto.City, dto.ZipCode, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Locker{Id: id}
}

func DeleteLocker(id string) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %s", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

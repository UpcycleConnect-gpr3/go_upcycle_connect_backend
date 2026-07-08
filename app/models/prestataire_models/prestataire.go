package prestataire_models

import (
	"fmt"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"
)

const TABLE = "PRESTATAIRES"

type Prestataire struct {
	Id        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Type      string `db:"type" json:"type"`
	Email     string `db:"email" json:"email"`
	Phone     string `db:"phone" json:"phone"`
	City      string `db:"city" json:"city"`
	Status    string `db:"status" json:"status"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type CreatePrestataireDTO struct {
	Name   string
	Type   string
	Email  string
	Phone  string
	City   string
	Status string
}

type UpdatePrestataireDTO struct {
	Name   string
	Type   string
	Email  string
	Phone  string
	City   string
	Status string
}

func (prestataire *Prestataire) Get(columns []string, by string, value ...any) error {
	return db.GetQuery[Prestataire](database.UpcycleConnect, TABLE, prestataire, columns, by, value...)
}

func (prestataire *Prestataire) All(columns []string, dest *[]Prestataire) error {
	return db.AllQuery[Prestataire](database.UpcycleConnect, TABLE, columns, dest)
}

func CreatePrestataire(dto CreatePrestataireDTO) *Prestataire {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Name)
	status := dto.Status
	if status == "" {
		status = "active"
	}
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (name, type, email, phone, city, status) VALUES (?, ?, ?, ?, ?, ?)",
		dto.Name, dto.Type, dto.Email, dto.Phone, dto.City, status,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &Prestataire{Id: int(id)}
}

func UpdatePrestataire(id int, dto UpdatePrestataireDTO) *Prestataire {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET name=?, type=?, email=?, phone=?, city=?, status=? WHERE id=?",
		dto.Name, dto.Type, dto.Email, dto.Phone, dto.City, dto.Status, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Prestataire{Id: int(id)}
}

func DeletePrestataire(id int) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

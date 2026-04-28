package order_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"

	"github.com/google/uuid"
)

const TABLE = "ORDERS"

type Order struct {
	Id        string `db:"id" json:"id"`
	Street    string `db:"street" json:"street"`
	City      string `db:"city" json:"city"`
	ZipCode   string `db:"zip_code" json:"zip_code"`
	UserId    string `db:"user_id" json:"user_id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type CreateOrderDTO struct {
	Street  string
	City    string
	ZipCode string
	UserId  string
}

type UpdateOrderDTO struct {
	Street  string
	City    string
	ZipCode string
	UserId  string
}

func (m *Order) Get(columns []string, by string, value any) error {
	return db.GetQuery[Order](database.UpcycleConnect, TABLE, columns, by, value, m)
}

func (m *Order) All(columns []string, dest *[]Order) error {
	return db.AllQuery[Order](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateOrder(dto CreateOrderDTO) *Order {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Street)
	orderId := uuid.New().String()
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (id, street, city, zip_code, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())",
		orderId, dto.Street, dto.City, dto.ZipCode, dto.UserId,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Order{Id: orderId}
}

func UpdateOrder(id string, dto UpdateOrderDTO) *Order {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %s", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET street=?, city=?, zip_code=?, user_id=?, updated_at=NOW() WHERE id=?",
		dto.Street, dto.City, dto.ZipCode, dto.UserId, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Order{Id: id}
}

func DeleteOrder(id string) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %s", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

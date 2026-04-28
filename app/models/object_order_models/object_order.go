package object_order_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
)

const TABLE = "OBJECT_ORDER"

type ObjectOrder struct {
	Id        int    `db:"id" json:"id"`
	ObjectId  string `db:"object_id" json:"object_id"`
	OrderId   string `db:"order_id" json:"order_id"`
	Amount    int    `db:"amount" json:"amount"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type CreateObjectOrderDTO struct {
	ObjectId string
	OrderId  string
	Amount   int
}

type UpdateObjectOrderDTO struct {
	Amount int
}

func (m *ObjectOrder) Get(columns []string, by string, value any) error {
	return db.GetQuery[ObjectOrder](database.UpcycleConnect, TABLE, columns, by, value, m)
}

func (m *ObjectOrder) All(columns []string, dest *[]ObjectOrder) error {
	return db.AllQuery[ObjectOrder](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateObjectOrder(dto CreateObjectOrderDTO) *ObjectOrder {
	action := fmt.Sprintf("INSERT INTO %s: object_id=%s, order_id=%s", TABLE, dto.ObjectId, dto.OrderId)
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (object_id, order_id, amount, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())",
		dto.ObjectId, dto.OrderId, dto.Amount,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &ObjectOrder{Id: int(id)}
}

func DeleteObjectOrder(id int) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

func GetObjectOrderByObjectAndOrder(objectId string, orderId string) *ObjectOrder {
	var oo ObjectOrder
	if err := oo.Get([]string{"id", "object_id", "order_id", "amount"}, "object_id=? AND order_id=?", []any{objectId, orderId}); err != nil {
		return nil
	}
	return &oo
}

func GetObjectOrdersByOrder(orderId string) []ObjectOrder {
	result := []ObjectOrder{}
	_ = database.UpcycleConnect.Select(&result, "SELECT * FROM "+TABLE+" WHERE order_id=?", orderId)
	return result
}

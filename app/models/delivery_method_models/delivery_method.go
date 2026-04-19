package delivery_method_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
)

const TABLE = "DELIVERY_METHODS"

type DeliveryMethod struct {
	Id        int     `db:"id" json:"id"`
	Name      string  `db:"name" json:"name"`
	Cost      float64 `db:"cost" json:"cost"`
	CreatedAt string  `db:"created_at" json:"created_at"`
	UpdatedAt string  `db:"updated_at" json:"updated_at"`
}

type CreateDeliveryMethodDTO struct {
	Name string
	Cost float64
}

type UpdateDeliveryMethodDTO struct {
	Name string
	Cost float64
}

type ObjectSummary struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (deliveryMethod *DeliveryMethod) Get(columns []string, by string, value any) error {
	return db.GetQuery[DeliveryMethod](database.UpcycleConnect, TABLE, columns, by, value, deliveryMethod)
}

func (deliveryMethod *DeliveryMethod) All(columns []string, dest *[]DeliveryMethod) error {
	return db.AllQuery[DeliveryMethod](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateDeliveryMethod(dto CreateDeliveryMethodDTO) *DeliveryMethod {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Name)
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (name, cost) VALUES (?, ?)",
		dto.Name, dto.Cost,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &DeliveryMethod{Id: int(id)}
}

func UpdateDeliveryMethod(id int, dto UpdateDeliveryMethodDTO) *DeliveryMethod {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET name=?, cost=? WHERE id=?",
		dto.Name, dto.Cost, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &DeliveryMethod{Id: int(id)}
}

func DeleteDeliveryMethod(id int) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

func GetDeliveryMethodObjects(deliveryMethodID int) []ObjectSummary {
	result := []ObjectSummary{}
	rows, err := database.UpcycleConnect.Query(
		"SELECT o.id, o.name FROM OBJECTS o JOIN OBJECT_DELIVERY_METHOD odm ON o.id=odm.object_id WHERE odm.delivery_method_id=?",
		deliveryMethodID,
	)
	if err != nil {
		log.Database("SELECT DELIVERY METHOD OBJECTS", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		o := ObjectSummary{}
		_ = rows.Scan(&o.Id, &o.Name)
		result = append(result, o)
	}
	return result
}

package order_delivery_method_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
)

const TABLE = "ORDER_DELIVERY_METHOD"

type OrderDeliveryMethod struct {
	OrderId          string  `db:"order_id" json:"order_id"`
	DeliveryMethodId int     `db:"delivery_method_id" json:"delivery_method_id"`
	Price            float64 `db:"price" json:"price"`
	CreatedAt        string  `db:"created_at" json:"created_at"`
	UpdatedAt        string  `db:"updated_at" json:"updated_at"`
}

type CreateOrderDeliveryMethodDTO struct {
	OrderId          string
	DeliveryMethodId int
	Price            float64
}

type UpdateOrderDeliveryMethodDTO struct {
	Price float64
}

func (m *OrderDeliveryMethod) Get(columns []string, by string, value any) error {
	return db.GetQuery[OrderDeliveryMethod](database.UpcycleConnect, TABLE, columns, by, value, m)
}

func (m *OrderDeliveryMethod) All(columns []string, dest *[]OrderDeliveryMethod) error {
	return db.AllQuery[OrderDeliveryMethod](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateOrderDeliveryMethod(dto CreateOrderDeliveryMethodDTO) *OrderDeliveryMethod {
	action := fmt.Sprintf("INSERT INTO %s: order_id=%s, delivery_method_id=%d", TABLE, dto.OrderId, dto.DeliveryMethodId)
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (order_id, delivery_method_id, price, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())",
		dto.OrderId, dto.DeliveryMethodId, dto.Price,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &OrderDeliveryMethod{OrderId: dto.OrderId, DeliveryMethodId: dto.DeliveryMethodId, Price: dto.Price}
}

func DeleteOrderDeliveryMethod(orderId string, deliveryMethodId int) {
	action := fmt.Sprintf("DELETE FROM %s WHERE order_id=%s AND delivery_method_id=%d", TABLE, orderId, deliveryMethodId)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE order_id=? AND delivery_method_id=?", orderId, deliveryMethodId)
	if err != nil {
		log.Database(action, err)
	}
}

func GetOrderDeliveryMethodsByOrder(orderId string) []OrderDeliveryMethod {
	result := []OrderDeliveryMethod{}
	_ = database.UpcycleConnect.Select(&result, "SELECT * FROM "+TABLE+" WHERE order_id=?", orderId)
	return result
}

func GetOrderDeliveryMethod(orderId string, deliveryMethodId int) *OrderDeliveryMethod {
	var odm OrderDeliveryMethod
	if err := odm.Get([]string{"order_id", "delivery_method_id", "price"}, "order_id=? AND delivery_method_id=?", []any{orderId, deliveryMethodId}); err != nil {
		return nil
	}
	return &odm
}

func UpdateOrderDeliveryMethod(orderId string, deliveryMethodId int, dto UpdateOrderDeliveryMethodDTO) *OrderDeliveryMethod {
	action := fmt.Sprintf("UPDATE %s WHERE order_id=%s AND delivery_method_id=%d", TABLE, orderId, deliveryMethodId)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET price=?, updated_at=NOW() WHERE order_id=? AND delivery_method_id=?",
		dto.Price, orderId, deliveryMethodId,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &OrderDeliveryMethod{OrderId: orderId, DeliveryMethodId: deliveryMethodId, Price: dto.Price}
}

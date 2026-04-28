package order_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/order_actions"
	"go-upcycle_connect-backend/app/models/order_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func IndexOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var order order_models.Order
	var orders []order_models.Order
	columns := []string{"id", "street", "city", "zip_code", "user_id", "created_at", "updated_at"}
	if err := order.All(columns, &orders); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, orders)
}

func ShowOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrOrderNotFound, http.StatusNotFound)
		return
	}
	var order order_models.Order
	columns := []string{"id", "street", "city", "zip_code", "user_id", "created_at", "updated_at"}
	if err := order.Get(columns, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrOrderNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, order)
}

func StoreOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto order_models.CreateOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, order := order_actions.CreateOrder(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if order == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]string{"id": order.Id})
}

func UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrOrderNotFound, http.StatusNotFound)
		return
	}
	var order order_models.Order
	if err := order.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrOrderNotFound, http.StatusNotFound)
		return
	}
	var dto order_models.UpdateOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, updated := order_actions.UpdateOrder(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrOrderNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]string{"id": updated.Id})
}

func DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrOrderNotFound, http.StatusNotFound)
		return
	}
	var order order_models.Order
	if err := order.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrOrderNotFound, http.StatusNotFound)
		return
	}
	order_models.DeleteOrder(id)
	response.NewSuccessMessage(w, "Order deleted successfully")
}

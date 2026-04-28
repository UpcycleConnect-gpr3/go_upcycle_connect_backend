package object_order_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/object_order_actions"
	"go-upcycle_connect-backend/app/models/object_order_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func IndexObjectOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var oo object_order_models.ObjectOrder
	var objectOrders []object_order_models.ObjectOrder
	columns := []string{"id", "object_id", "order_id", "amount", "created_at", "updated_at"}
	if err := oo.All(columns, &objectOrders); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, objectOrders)
}

func ShowObjectOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var oo object_order_models.ObjectOrder
	columns := []string{"id", "object_id", "order_id", "amount", "created_at", "updated_at"}
	if err := oo.Get(columns, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrObjectOrderNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, oo)
}

func StoreObjectOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto object_order_models.CreateObjectOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, objectOrder := object_order_actions.CreateObjectOrder(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if objectOrder == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": objectOrder.Id})
}

func DeleteObjectOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var oo object_order_models.ObjectOrder
	if err := oo.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrObjectOrderNotFound, http.StatusNotFound)
		return
	}
	object_order_models.DeleteObjectOrder(id)
	response.NewSuccessMessage(w, "Object-Order deleted successfully")
}

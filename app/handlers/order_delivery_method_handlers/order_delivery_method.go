package order_delivery_method_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/order_delivery_method_actions"
	"go-upcycle_connect-backend/app/models/order_delivery_method_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func IndexOrderDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var odm order_delivery_method_models.OrderDeliveryMethod
	var orderDeliveryMethods []order_delivery_method_models.OrderDeliveryMethod
	columns := []string{"order_id", "delivery_method_id", "price", "created_at", "updated_at"}
	if err := odm.All(columns, &orderDeliveryMethods); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, orderDeliveryMethods)
}

func StoreOrderDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto order_delivery_method_models.CreateOrderDeliveryMethodDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, odm := order_delivery_method_actions.CreateOrderDeliveryMethod(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if odm == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, odm)
}

func DeleteOrderDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	orderId := request.Request(r, "orderId").Value()
	if orderId == "" {
		response.NewErrorMessage(w, response.ErrOrderNotFound, http.StatusNotFound)
		return
	}
	deliveryMethodId := request.Request(r, "deliveryMethodId").ConvertToInt(w)
	if deliveryMethodId == -1 {
		return
	}
	order_delivery_method_models.DeleteOrderDeliveryMethod(orderId, deliveryMethodId)
	response.NewSuccessMessage(w, "Order-DeliveryMethod deleted successfully")
}

package delivery_method_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/delivery_method_actions"
	"go-upcycle_connect-backend/app/models/delivery_method_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func GetDeliveryMethodsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dm delivery_method_models.DeliveryMethod
	var deliveryMethods []delivery_method_models.DeliveryMethod

	columns := []string{"id", "name", "cost", "created_at", "updated_at"}
	if err := dm.All(columns, &deliveryMethods); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, deliveryMethods)
}

func GetDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var dm delivery_method_models.DeliveryMethod
	columns := []string{"id", "name", "cost", "created_at", "updated_at"}
	if err := dm.Get(columns, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, dm)
}

func CreateDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto delivery_method_actions.CreateDeliveryMethodDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, deliveryMethod := delivery_method_actions.CreateDeliveryMethod(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if deliveryMethod == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": deliveryMethod.Id})
}

func UpdateDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var dm delivery_method_models.DeliveryMethod
	if err := dm.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusNotFound)
		return
	}
	var dto delivery_method_actions.UpdateDeliveryMethodDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, updated := delivery_method_actions.UpdateDeliveryMethod(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": updated.Id})
}

func DeleteDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var dm delivery_method_models.DeliveryMethod
	if err := dm.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusNotFound)
		return
	}
	delivery_method_models.DeleteDeliveryMethod(id)
	response.NewSuccessMessage(w, "Delivery method deleted")
}

func GetDeliveryMethodObjectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var dm delivery_method_models.DeliveryMethod
	if err := dm.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusNotFound)
		return
	}
	objects := delivery_method_models.GetDeliveryMethodObjects(id)
	response.NewSuccessData(w, objects)
}

package delivery_method_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/delivery_method_actions"
	"go-upcycle_connect-backend/app/models/delivery_method_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
	"strconv"
)

func GetDeliveryMethodsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dm delivery_method_models.DeliveryMethod
	var deliveryMethods []delivery_method_models.DeliveryMethod

	err := dm.All([]string{"id", "name", "cost", "created_at", "updated_at"}, &deliveryMethods)
	if err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, deliveryMethods)
}

func GetDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusBadRequest)
		return
	}
	var dm delivery_method_models.DeliveryMethod
	err = dm.Get([]string{"id", "name", "cost", "created_at", "updated_at"}, "id", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, dm)
}

func CreateDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto delivery_method_actions.CreateDeliveryMethodDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	dm, errs := delivery_method_actions.CreateDeliveryMethod(dto)
	if len(errs) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, errs)
		return
	}
	if dm == nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, dm)
}

func UpdateDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusBadRequest)
		return
	}
	var dm delivery_method_models.DeliveryMethod
	err = dm.Get([]string{"id"}, "id", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusNotFound)
		return
	}
	var dto delivery_method_actions.UpdateDeliveryMethodDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	updated, errs := delivery_method_actions.UpdateDeliveryMethod(id, dto)
	if len(errs) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, errs)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, updated)
}

func DeleteDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusBadRequest)
		return
	}
	var dm delivery_method_models.DeliveryMethod
	err = dm.Get([]string{"id"}, "id", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusNotFound)
		return
	}
	delivery_method_models.DeleteDeliveryMethod(id)
	response.NewSuccessMessage(w, response.SuccessDeleted)
}

func GetDeliveryMethodObjectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusBadRequest)
		return
	}
	var dm delivery_method_models.DeliveryMethod
	err = dm.Get([]string{"id"}, "id", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrDeliveryMethodNotFound, http.StatusNotFound)
		return
	}
	objects := delivery_method_models.GetDeliveryMethodObjects(id)
	response.NewSuccessData(w, objects)
}

package object_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/object_actions"
	"go-upcycle_connect-backend/app/models/object_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"

	"github.com/google/uuid"
)

func findObject(w http.ResponseWriter, id int) bool {
	var obj object_models.Object
	if err := obj.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return false
	}
	return true
}

func GetObjectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var obj object_models.Object
	var objects []object_models.Object
	columns := []string{"id", "name", "material", "`condition`", "description", "upcycling_score", "created_at", "updated_at"}
	if err := obj.All(columns, &objects); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, objects)
}

func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var obj object_models.Object
	columns := []string{"id", "name", "material", "`condition`", "description", "upcycling_score", "created_at", "updated_at"}
	if err := obj.Get(columns, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, obj)
}

func CreateObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto object_actions.CreateObjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, obj := object_actions.CreateObject(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if obj == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": obj.Id})
}

func UpdateObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findObject(w, id) {
		return
	}
	var dto object_actions.UpdateObjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, updated := object_actions.UpdateObject(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": updated.Id})
}

func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findObject(w, id) {
		return
	}
	object_models.DeleteObject(id)
	response.NewSuccessMessage(w, "Object deleted")
}

func GetObjectScoreHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	score := object_models.GetObjectScore(id)
	if score == nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, score)
}

func GetObjectDeliveryMethodsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findObject(w, id) {
		return
	}
	dms := object_models.GetObjectDeliveryMethods(id)
	response.NewSuccessData(w, dms)
}

func LinkDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	dmId := request.Request(r, "deliveryMethodId").ConvertToInt(w)
	if dmId == -1 {
		return
	}
	object_actions.LinkDeliveryMethod(id, dmId)
	response.NewSuccessMessage(w, response.SuccessLinked)
}

func UnlinkDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	dmId := request.Request(r, "deliveryMethodId").ConvertToInt(w)
	if dmId == -1 {
		return
	}
	object_actions.UnlinkDeliveryMethod(id, dmId)
	response.NewSuccessMessage(w, response.SuccessUnlinked)
}

func GetObjectProjectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findObject(w, id) {
		return
	}
	projects := object_models.GetObjectProjects(id)
	response.NewSuccessData(w, projects)
}

func LinkProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	projectId := request.Request(r, "projectId").ConvertToInt(w)
	if projectId == -1 {
		return
	}
	object_actions.LinkProject(id, projectId)
	response.NewSuccessMessage(w, response.SuccessLinked)
}

func UnlinkProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	projectId := request.Request(r, "projectId").ConvertToInt(w)
	if projectId == -1 {
		return
	}
	object_actions.UnlinkProject(id, projectId)
	response.NewSuccessMessage(w, response.SuccessUnlinked)
}

func GetObjectUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findObject(w, id) {
		return
	}
	users := object_models.GetObjectUsers(id)
	response.NewSuccessData(w, users)
}

func LinkUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	userId := request.Request(r, "userId").Value()
	if userId == "" {
		response.NewErrorMessage(w, response.ErrUserNotFound, http.StatusBadRequest)
		return
	}
	if _, err := uuid.Parse(userId); err != nil {
		response.NewErrorMessage(w, "Invalid user ID format (expected UUID)", http.StatusBadRequest)
		return
	}
	if err := object_actions.LinkUser(id, userId); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessMessage(w, response.SuccessLinked)
}

func UnlinkUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	userId := request.Request(r, "userId").Value()
	if userId == "" {
		response.NewErrorMessage(w, response.ErrUserNotFound, http.StatusBadRequest)
		return
	}
	if _, err := uuid.Parse(userId); err != nil {
		response.NewErrorMessage(w, "Invalid user ID format (expected UUID)", http.StatusBadRequest)
		return
	}
	if err := object_actions.UnlinkUser(id, userId); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessMessage(w, response.SuccessUnlinked)
}

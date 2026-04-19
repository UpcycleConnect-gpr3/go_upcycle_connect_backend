package object_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/object_actions"
	"go-upcycle_connect-backend/app/models/object_models"
	"go-upcycle_connect-backend/utils/jwt"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
	"strconv"
)

func parseObjectID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func findObject(w http.ResponseWriter, id int) bool {
	var obj object_models.Object
	if err := obj.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return false
	}
	return true
}

func parsePathInt(w http.ResponseWriter, r *http.Request, key string, notFoundErr string) (int, bool) {
	id, err := strconv.Atoi(r.PathValue(key))
	if err != nil {
		response.NewErrorMessage(w, notFoundErr, http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func GetObjectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	var obj object_models.Object
	var objects []object_models.Object
	if err := obj.All([]string{"id", "name", "material", "`condition`", "description", "upcycling_score", "created_at", "updated_at"}, &objects); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, objects)
}

func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
		return
	}
	var obj object_models.Object
	if err := obj.Get([]string{"id", "name", "material", "`condition`", "description", "upcycling_score", "created_at", "updated_at"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, obj)
}

func CreateObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	var dto object_actions.CreateObjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	obj, errs := object_actions.CreateObject(dto)
	if len(errs) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, errs)
		return
	}
	if obj == nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, obj)
}

func UpdateObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
		return
	}
	if !findObject(w, id) {
		return
	}
	var dto object_actions.UpdateObjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	updated, errs := object_actions.UpdateObject(id, dto)
	if len(errs) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, errs)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, updated)
}

func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
		return
	}
	if !findObject(w, id) {
		return
	}
	object_models.DeleteObject(id)
	response.NewSuccessMessage(w, response.SuccessDeleted)
}

func GetObjectScoreHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
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
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
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
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
		return
	}
	dmId, ok := parsePathInt(w, r, "deliveryMethodId", response.ErrDeliveryMethodNotFound)
	if !ok {
		return
	}
	object_actions.LinkDeliveryMethod(id, dmId)
	response.NewSuccessMessage(w, response.SuccessLinked)
}

func UnlinkDeliveryMethodHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
		return
	}
	dmId, ok := parsePathInt(w, r, "deliveryMethodId", response.ErrDeliveryMethodNotFound)
	if !ok {
		return
	}
	object_actions.UnlinkDeliveryMethod(id, dmId)
	response.NewSuccessMessage(w, response.SuccessUnlinked)
}

func GetObjectProjectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
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
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
		return
	}
	projectId, ok := parsePathInt(w, r, "projectId", response.ErrProjectNotFound)
	if !ok {
		return
	}
	object_actions.LinkProject(id, projectId)
	response.NewSuccessMessage(w, response.SuccessLinked)
}

func UnlinkProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
		return
	}
	projectId, ok := parsePathInt(w, r, "projectId", response.ErrProjectNotFound)
	if !ok {
		return
	}
	object_actions.UnlinkProject(id, projectId)
	response.NewSuccessMessage(w, response.SuccessUnlinked)
}

func GetObjectUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
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
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
		return
	}
	userId := r.PathValue("userId")
	if userId == "" {
		response.NewErrorMessage(w, response.ErrUserNotFound, http.StatusBadRequest)
		return
	}
	object_actions.LinkUser(id, userId)
	response.NewSuccessMessage(w, response.SuccessLinked)
}

func UnlinkUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if !jwt.Auth(w, r) {
		return
	}
	id, ok := parseObjectID(w, r)
	if !ok {
		return
	}
	userId := r.PathValue("userId")
	if userId == "" {
		response.NewErrorMessage(w, response.ErrUserNotFound, http.StatusBadRequest)
		return
	}
	object_actions.UnlinkUser(id, userId)
	response.NewSuccessMessage(w, response.SuccessUnlinked)
}

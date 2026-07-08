package object_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/object_actions"
	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/object_models"
	"go-upcycle_connect-backend/app/models/score_models"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"go-upcycle_connect-backend/utils/rules"
	"net/http"

	"github.com/google/uuid"
)

func findObject(w http.ResponseWriter, id string) bool {
	var obj object_models.Object
	if err := obj.Get([]string{"id"}, db.IdClause, id); err != nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return false
	}
	return true
}

func IndexObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var obj object_models.Object
	var objects []object_models.Object
	columns := []string{"id", "name", "description", "price", "image_path", "column_for_calc_the_score", "category", "item_condition", "quantity", "user_id", "score", "created_at", "updated_at"}
	if err := obj.All(columns, &objects); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, objects)
}

func ShowObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return
	}
	var obj object_models.Object
	columns := []string{"id", "name", "description", "price", "image_path", "column_for_calc_the_score", "category", "item_condition", "quantity", "user_id", "score", "created_at", "updated_at"}
	if err := obj.Get(columns, db.IdClause, id); err != nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, obj)
}

func StoreObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto object_models.CreateObjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	// L'auteur de l'annonce est l'utilisateur du token, jamais le body.
	dto.UserId = auth_middleware.GetUserId(r.Context())

	// Categorie/etat : defauts + validation, puis Upcycler Score calcule.
	if dto.Category == "" {
		dto.Category = score_models.DefaultCategory
	}
	if dto.Condition == "" {
		dto.Condition = score_models.DefaultCondition
	}
	if !score_models.IsValidCategory(dto.Category) {
		response.NewValidationError(w, response.ErrInvalidBody, []rules.ValidationError{{Field: "category", Message: "invalid category"}})
		return
	}
	if !score_models.IsValidCondition(dto.Condition) {
		response.NewValidationError(w, response.ErrInvalidBody, []rules.ValidationError{{Field: "condition", Message: "invalid condition"}})
		return
	}
	dto.Score = score_models.Compute(dto.Category, dto.Condition)

	validationErrors, obj := object_actions.CreateObject(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if obj == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]string{"id": obj.Id})
}

func UpdateObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return
	}
	if !findObject(w, id) {
		return
	}
	var dto object_models.UpdateObjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	// Reprend la categorie/etat existants si absents du body, valide, puis
	// recalcule l'Upcycler Score.
	var current object_models.Object
	if err := current.Get([]string{"category", "item_condition"}, db.IdClause, id); err == nil {
		if dto.Category == "" {
			dto.Category = current.Category
		}
		if dto.Condition == "" {
			dto.Condition = current.Condition
		}
	}
	if dto.Category == "" {
		dto.Category = score_models.DefaultCategory
	}
	if dto.Condition == "" {
		dto.Condition = score_models.DefaultCondition
	}
	if !score_models.IsValidCategory(dto.Category) {
		response.NewValidationError(w, response.ErrInvalidBody, []rules.ValidationError{{Field: "category", Message: "invalid category"}})
		return
	}
	if !score_models.IsValidCondition(dto.Condition) {
		response.NewValidationError(w, response.ErrInvalidBody, []rules.ValidationError{{Field: "condition", Message: "invalid condition"}})
		return
	}
	dto.Score = score_models.Compute(dto.Category, dto.Condition)

	validationErrors, updated := object_actions.UpdateObject(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]string{"id": updated.Id})
}

func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
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
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
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
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
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
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
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
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
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
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
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
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
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
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
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
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
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
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
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
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
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

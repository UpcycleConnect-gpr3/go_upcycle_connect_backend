package project_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/project_actions"
	"go-upcycle_connect-backend/app/models/project_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func findProject(w http.ResponseWriter, id int) bool {
	var p project_models.Project
	if err := p.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrProjectNotFound, http.StatusNotFound)
		return false
	}
	return true
}

func IndexProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var p project_models.Project
	var projects []project_models.Project
	columns := []string{"id", "name", "description", "image_path", "user_id", "created_at", "updated_at"}
	if err := p.All(columns, &projects); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, projects)
}

func ShowProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var p project_models.Project
	columns := []string{"id", "name", "description", "image_path", "user_id", "created_at", "updated_at"}
	if err := p.Get(columns, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrProjectNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, p)
}

func StoreProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto project_models.CreateProjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, project := project_actions.CreateProject(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if project == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": project.Id})
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findProject(w, id) {
		return
	}
	var dto project_models.UpdateProjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, updated := project_actions.UpdateProject(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrProjectNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": updated.Id})
}

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findProject(w, id) {
		return
	}
	project_models.DeleteProject(id)
	response.NewSuccessMessage(w, "Project deleted")
}

func GetProjectObjectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findProject(w, id) {
		return
	}
	objects := project_models.GetProjectObjects(id)
	response.NewSuccessData(w, objects)
}

func LinkObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	objectId := request.Request(r, "objectId").Value()
	if objectId == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusBadRequest)
		return
	}
	userId := request.Request(r, "userId").Value()
	if userId == "" {
		userId = "" // Allow empty userId for now
	}
	project_actions.LinkObject(id, objectId, userId)
	response.NewSuccessMessage(w, response.SuccessLinked)
}

func UnlinkObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	objectId := request.Request(r, "objectId").Value()
	if objectId == "" {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusBadRequest)
		return
	}
	project_models.UnlinkObject(id, objectId)
	response.NewSuccessMessage(w, response.SuccessUnlinked)
}

func GetProjectStepsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findProject(w, id) {
		return
	}
	steps := project_models.GetProjectSteps(id)
	response.NewSuccessData(w, steps)
}

func CreateProjectStepHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if !findProject(w, id) {
		return
	}
	var dto struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ImagePath   string `json:"image_path"`
		ScheduledAt string `json:"scheduled_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, projectStep := project_actions.CreateProjectStep(id, dto.Name, dto.Description, dto.ImagePath, dto.ScheduledAt)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if projectStep == nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, projectStep)
}

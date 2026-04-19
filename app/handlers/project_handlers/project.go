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

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var p project_models.Project
	var projects []project_models.Project
	columns := []string{"id", "name", "description", "score", "created_at", "updated_at"}
	if err := p.All(columns, &projects); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, projects)
}

func GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var p project_models.Project
	columns := []string{"id", "name", "description", "score", "created_at", "updated_at"}
	if err := p.Get(columns, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrProjectNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, p)
}

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto project_actions.CreateProjectDTO
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
	var dto project_actions.UpdateProjectDTO
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

func GetProjectScoreHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	score := project_models.GetProjectScore(id)
	if score == nil {
		response.NewErrorMessage(w, response.ErrProjectNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, score)
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
	objectId := request.Request(r, "objectId").ConvertToInt(w)
	if objectId == -1 {
		return
	}
	project_actions.LinkObject(id, objectId)
	response.NewSuccessMessage(w, response.SuccessLinked)
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
	var dto project_actions.CreateProjectStepDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, projectStep := project_actions.CreateProjectStep(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if projectStep == nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": projectStep.Id})
}

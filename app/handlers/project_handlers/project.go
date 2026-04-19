package project_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/project_actions"
	"go-upcycle_connect-backend/app/models/project_models"
	"go-upcycle_connect-backend/utils/handler"
	"go-upcycle_connect-backend/utils/log"
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
	if err := p.All([]string{"id", "name", "description", "score", "created_at", "updated_at"}, &projects); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, projects)
}

func GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id, ok := handler.ParsePathInt(w, r, "id", response.ErrProjectNotFound)
	if !ok {
		return
	}
	var p project_models.Project
	if err := p.Get([]string{"id", "name", "description", "score", "created_at", "updated_at"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrProjectNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, p)
}

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto project_actions.CreateProjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	p, errs := project_actions.CreateProject(dto)
	if len(errs) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, errs)
		return
	}
	if p == nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, p)
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id, ok := handler.ParsePathInt(w, r, "id", response.ErrProjectNotFound)
	if !ok {
		return
	}
	if !findProject(w, id) {
		return
	}
	var dto project_actions.UpdateProjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	updated, errs := project_actions.UpdateProject(id, dto)
	if len(errs) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, errs)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrProjectNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, updated)
}

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id, ok := handler.ParsePathInt(w, r, "id", response.ErrProjectNotFound)
	if !ok {
		return
	}
	if !findProject(w, id) {
		return
	}
	project_models.DeleteProject(id)
	response.NewSuccessMessage(w, response.SuccessDeleted)
}

func GetProjectScoreHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id, ok := handler.ParsePathInt(w, r, "id", response.ErrProjectNotFound)
	if !ok {
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
	id, ok := handler.ParsePathInt(w, r, "id", response.ErrProjectNotFound)
	if !ok {
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
	id, ok := handler.ParsePathInt(w, r, "id", response.ErrProjectNotFound)
	if !ok {
		return
	}
	objectId, ok := handler.ParsePathInt(w, r, "objectId", response.ErrObjectNotFound)
	if !ok {
		return
	}
	project_actions.LinkObject(id, objectId)
	response.NewSuccessMessage(w, response.SuccessLinked)
}

func GetProjectStepsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id, ok := handler.ParsePathInt(w, r, "id", response.ErrProjectNotFound)
	if !ok {
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
	id, ok := handler.ParsePathInt(w, r, "id", response.ErrProjectNotFound)
	if !ok {
		return
	}
	if !findProject(w, id) {
		return
	}
	var dto project_actions.CreateProjectStepDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	s, errs := project_actions.CreateProjectStep(id, dto)
	if len(errs) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, errs)
		return
	}
	if s == nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, s)
}

package prestataire_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/prestataire_actions"
	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/activity_models"
	"go-upcycle_connect-backend/app/models/prestataire_models"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
	"strconv"
)

func IndexPrestataireHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var p prestataire_models.Prestataire
	var prestataires []prestataire_models.Prestataire

	columns := []string{"id", "name", "type", "email", "phone", "city", "status", "created_at", "updated_at"}
	if err := p.All(columns, &prestataires); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, prestataires)
}

func ShowPrestataireHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var p prestataire_models.Prestataire
	columns := []string{"id", "name", "type", "email", "phone", "city", "status", "created_at", "updated_at"}
	if err := p.Get(columns, db.IdClause, id); err != nil {
		response.NewErrorMessage(w, response.ErrPrestataireNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, p)
}

func StorePrestataireHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto prestataire_models.CreatePrestataireDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, prestataire := prestataire_actions.CreatePrestataire(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if prestataire == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	activity_models.Record(auth_middleware.GetUserId(r.Context()), "create", "prestataire", strconv.Itoa(prestataire.Id), dto.Name)
	response.NewSuccessData(w, map[string]int{"id": prestataire.Id})
}

func UpdatePrestataireHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var p prestataire_models.Prestataire
	if err := p.Get([]string{"id"}, db.IdClause, id); err != nil {
		response.NewErrorMessage(w, response.ErrPrestataireNotFound, http.StatusNotFound)
		return
	}
	var dto prestataire_models.UpdatePrestataireDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, updated := prestataire_actions.UpdatePrestataire(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrPrestataireNotFound, http.StatusInternalServerError)
		return
	}
	activity_models.Record(auth_middleware.GetUserId(r.Context()), "update", "prestataire", strconv.Itoa(updated.Id), dto.Name)
	response.NewSuccessData(w, map[string]int{"id": updated.Id})
}

func DeletePrestataireHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	var p prestataire_models.Prestataire
	if err := p.Get([]string{"id"}, db.IdClause, id); err != nil {
		response.NewErrorMessage(w, response.ErrPrestataireNotFound, http.StatusNotFound)
		return
	}
	prestataire_models.DeletePrestataire(id)
	activity_models.Record(auth_middleware.GetUserId(r.Context()), "delete", "prestataire", strconv.Itoa(id), "")
	response.NewSuccessMessage(w, "Prestataire deleted")
}

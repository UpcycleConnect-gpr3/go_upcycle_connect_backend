package package_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/package_actions"
	"go-upcycle_connect-backend/app/models/package_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func IndexPackageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var pkg package_models.Package
	var packages []package_models.Package
	columns := []string{"id", "weight", "code", "locker_id", "order_id", "created_at", "updated_at"}
	if err := pkg.All(columns, &packages); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, packages)
}

func ShowPackageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrPackageNotFound, http.StatusNotFound)
		return
	}
	var pkg package_models.Package
	columns := []string{"id", "weight", "code", "locker_id", "order_id", "created_at", "updated_at"}
	if err := pkg.Get(columns, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrPackageNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, pkg)
}

func StorePackageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto package_models.CreatePackageDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, pkg := package_actions.CreatePackage(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if pkg == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]string{"id": pkg.Id})
}

func UpdatePackageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrPackageNotFound, http.StatusNotFound)
		return
	}
	var pkg package_models.Package
	if err := pkg.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrPackageNotFound, http.StatusNotFound)
		return
	}
	var dto package_models.UpdatePackageDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, updated := package_actions.UpdatePackage(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrPackageNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]string{"id": updated.Id})
}

func DeletePackageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrPackageNotFound, http.StatusNotFound)
		return
	}
	var pkg package_models.Package
	if err := pkg.Get([]string{"id"}, "id", id); err != nil {
		response.NewErrorMessage(w, response.ErrPackageNotFound, http.StatusNotFound)
		return
	}
	package_models.DeletePackage(id)
	response.NewSuccessMessage(w, "Package deleted successfully")
}

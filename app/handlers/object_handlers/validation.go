package object_handlers

import (
	"net/http"

	"go-upcycle_connect-backend/app/models/object_models"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/jwt"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
)

func isAdmin(r *http.Request) bool {
	return jwt.RoleFromToken(r.Header.Get("Authorization")) == "administrator"
}

func setObjectValidated(w http.ResponseWriter, r *http.Request, validated bool) {
	log.Api(r)
	if !isAdmin(r) {
		response.NewErrorMessage(w, response.ErrForbidden, http.StatusForbidden)
		return
	}
	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusNotFound)
		return
	}
	var object object_models.Object
	if err := object.Get([]string{"id"}, db.IdClause, id); err != nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return
	}
	if err := object_models.SetAdValidated(id, validated); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]any{"id": id, "is_ad_validated": validated})
}

func ValidateObjectHandler(w http.ResponseWriter, r *http.Request) {
	setObjectValidated(w, r, true)
}

func RejectObjectHandler(w http.ResponseWriter, r *http.Request) {
	setObjectValidated(w, r, false)
}

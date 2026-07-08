package activity_handlers

import (
	"net/http"

	"go-upcycle_connect-backend/app/models/activity_models"
	"go-upcycle_connect-backend/utils/jwt"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
)

// GetLogsHandler — GET /logs (administrator)
// Journal d'activite du back office.
func GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if jwt.RoleFromToken(r.Header.Get("Authorization")) != "administrator" {
		response.NewErrorMessage(w, response.ErrForbidden, http.StatusForbidden)
		return
	}
	response.NewSuccessData(w, activity_models.GetAll(100))
}

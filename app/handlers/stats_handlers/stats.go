package stats_handlers

import (
	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/stats_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

// GetMyStatsHandler - GET /stats/me : statistiques de l'utilisateur du token.
func GetMyStatsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())
	response.NewSuccessData(w, stats_models.GetUserStats(userId))
}

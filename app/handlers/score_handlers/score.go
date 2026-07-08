package score_handlers

import (
	"go-upcycle_connect-backend/app/models/score_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

// GetScoreConfigHandler - GET /score/config
// Expose les categories (poids CO2) et etats (coefficients) pour que le
// frontend construise le formulaire d'annonce.
func GetScoreConfigHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	response.NewSuccessData(w, score_models.Config())
}

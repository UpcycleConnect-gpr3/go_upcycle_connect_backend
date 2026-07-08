package ad_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/actions/ad_actions"
	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/ad_models"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func findOwnedAd(w http.ResponseWriter, r *http.Request) *ad_models.Ad {
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return nil
	}
	var ad ad_models.Ad
	if err := ad.Get([]string{"id", "user_id"}, db.IdClause, id); err != nil {
		response.NewErrorMessage(w, response.ErrAdNotFound, http.StatusNotFound)
		return nil
	}
	if ad.UserId != auth_middleware.GetUserId(r.Context()) {
		response.NewErrorMessage(w, response.ErrAdNotFound, http.StatusNotFound)
		return nil
	}
	return &ad
}

func IndexAdHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())
	response.NewSuccessData(w, ad_models.GetUserAds(userId))
}

func StoreAdHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	var dto ad_models.CreateAdDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	dto.UserId = auth_middleware.GetUserId(r.Context())
	validationErrors, ad := ad_actions.CreateAd(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if ad == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": ad.Id})
}

func UpdateAdHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	ad := findOwnedAd(w, r)
	if ad == nil {
		return
	}
	var dto ad_models.UpdateAdDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	validationErrors, updated := ad_actions.UpdateAd(ad.Id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if updated == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]int{"id": updated.Id})
}

func DeleteAdHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	ad := findOwnedAd(w, r)
	if ad == nil {
		return
	}
	ad_models.DeleteAd(ad.Id)
	response.NewSuccessMessage(w, "Ad deleted successfully")
}

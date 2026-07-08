package deposit_handlers

import (
	"encoding/json"
	"net/http"

	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/locker_models"
	"go-upcycle_connect-backend/app/models/object_models"
	"go-upcycle_connect-backend/app/models/package_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
)

func SellerDeliveriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())
	response.NewSuccessData(w, package_models.GetSellerDeliveries(userId))
}

func BuyerDeliveriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())
	response.NewSuccessData(w, package_models.GetBuyerDeliveries(userId))
}

type depositConfirmDTO struct {
	Code string `json:"code"`
}

func DepositConfirmHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())

	var dto depositConfirmDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	if dto.Code == "" {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}

	pkg := package_models.GetByAnyCode(dto.Code)

	if pkg == nil || pkg.Code != dto.Code || pkg.Status != "awaiting_deposit" {
		response.NewErrorMessage(w, response.ErrPackageNotFound, http.StatusNotFound)
		return
	}

	var object object_models.Object
	if err := object.Get([]string{"id", "user_id", "score"}, "id = ?", pkg.ObjectId); err != nil || object.UserId != userId {
		response.NewErrorMessage(w, response.ErrForbidden, http.StatusForbidden)
		return
	}

	_ = package_models.MarkDeposited(pkg.Id)
	_ = locker_models.AdjustSlots(pkg.LockerId, -1)
	_ = object_models.SetStatusAndScore(pkg.ObjectId, "in_locker", object.Score)

	response.NewSuccessData(w, map[string]string{"status": "deposited", "object_id": pkg.ObjectId})
}

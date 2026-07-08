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

// SellerDeliveriesHandler — GET /packages/sales (auth)
// Ventes de l'utilisateur à déposer en casier (code de dépôt fourni).
func SellerDeliveriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())
	response.NewSuccessData(w, package_models.GetSellerDeliveries(userId))
}

// BuyerDeliveriesHandler — GET /packages/purchases (auth)
// Achats de l'utilisateur à récupérer en casier (code de retrait fourni).
func BuyerDeliveriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())
	response.NewSuccessData(w, package_models.GetBuyerDeliveries(userId))
}

type depositConfirmDTO struct {
	Code string `json:"code"`
}

// DepositConfirmHandler — POST /packages/deposit-confirm (auth)
// Le vendeur saisit son code de dépôt pour ouvrir le casier et y déposer
// l'objet vendu : la livraison passe 'deposited', le casier est décrémenté.
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
	// Doit être le code de dépôt (pkg.Code), une livraison en attente de dépôt.
	if pkg == nil || pkg.Code != dto.Code || pkg.Status != "awaiting_deposit" {
		response.NewErrorMessage(w, response.ErrPackageNotFound, http.StatusNotFound)
		return
	}

	// Vérifie que l'objet appartient bien au vendeur.
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

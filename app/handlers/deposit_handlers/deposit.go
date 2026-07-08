package deposit_handlers

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"time"

	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/locker_models"
	"go-upcycle_connect-backend/app/models/object_models"
	"go-upcycle_connect-backend/app/models/package_models"
	"go-upcycle_connect-backend/app/models/score_models"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
)

const codeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const depositTTLDays = 7

func generateCode() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = codeAlphabet[int(b[i])%len(codeAlphabet)]
	}
	return string(b)
}

func AvailableLockersHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	city := r.URL.Query().Get("city")
	response.NewSuccessData(w, locker_models.GetAvailableLockers(city))
}

type depositDTO struct {
	ObjectId string `json:"object_id"`
	LockerId string `json:"locker_id"`
	Weight   int    `json:"weight"`
}

func DepositHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	userId := auth_middleware.GetUserId(r.Context())

	var dto depositDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	if dto.ObjectId == "" || dto.LockerId == "" {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}

	var object object_models.Object
	if err := object.Get([]string{"id", "user_id", "category", "item_condition"}, db.IdClause, dto.ObjectId); err != nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return
	}
	if object.UserId != userId {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return
	}

	var locker locker_models.Locker
	if err := locker.Get([]string{"id", "available_slots"}, db.IdClause, dto.LockerId); err != nil {
		response.NewErrorMessage(w, response.ErrLockerNotFound, http.StatusNotFound)
		return
	}
	if locker.AvailableSlots <= 0 {
		response.NewErrorMessage(w, response.ErrLockerFull, http.StatusConflict)
		return
	}

	code := generateCode()
	expiry := time.Now().Add(depositTTLDays * 24 * time.Hour).Format("2006-01-02 15:04:05")

	pkg := package_models.CreateDeposit(dto.ObjectId, dto.LockerId, code, dto.Weight, expiry)
	if pkg == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}

	_ = locker_models.AdjustSlots(dto.LockerId, -1)

	score := score_models.Compute(object.Category, object.Condition)
	_ = object_models.SetStatusAndScore(dto.ObjectId, "in_locker", score)

	response.NewSuccessData(w, map[string]interface{}{
		"code":        code,
		"expiry_date": expiry,
		"score":       score,
		"package_id":  pkg.Id,
	})
}

type retrieveDTO struct {
	PackageCode string `json:"package_code"`
}

func RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	userId := auth_middleware.GetUserId(r.Context())

	var dto retrieveDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	if dto.PackageCode == "" {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}

	pkg := package_models.GetByAnyCode(dto.PackageCode)
	if pkg == nil {
		response.NewErrorMessage(w, response.ErrPackageNotFound, http.StatusNotFound)
		return
	}
	if pkg.Status != "deposited" {
		response.NewErrorMessage(w, response.ErrPackageAlreadyRetrieved, http.StatusConflict)
		return
	}
	if isExpired(pkg.ExpiryDate) {
		response.NewErrorMessage(w, response.ErrPackageExpired, http.StatusGone)
		return
	}

	newOwner := userId
	if pkg.BuyerId != "" {

		newOwner = pkg.BuyerId
	}
	_ = package_models.MarkRetrieved(pkg.Id)
	_ = locker_models.AdjustSlots(pkg.LockerId, 1)
	_ = object_models.SetStatusAndOwner(pkg.ObjectId, "retrieved", newOwner)

	response.NewSuccessData(w, map[string]string{
		"object_id": pkg.ObjectId,
		"status":    "retrieved",
	})
}

func PackageByCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	code := request.Request(r, "code").Value()
	if code == "" {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusNotFound)
		return
	}
	pkg := package_models.GetByCode(code)
	if pkg == nil {
		response.NewErrorMessage(w, response.ErrPackageNotFound, http.StatusNotFound)
		return
	}
	response.NewSuccessData(w, pkg)
}

func isExpired(expiry string) bool {
	if expiry == "" {
		return false
	}
	t, err := time.Parse("2006-01-02 15:04:05", expiry)
	if err != nil {

		t, err = time.Parse(time.RFC3339, expiry)
		if err != nil {
			return false
		}
	}
	return time.Now().After(t)
}

func DepositedPackagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	response.NewSuccessData(w, package_models.GetDepositedPackages())
}

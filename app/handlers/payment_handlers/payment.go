package payment_handlers

import (
	"crypto/rand"
	"encoding/json"
	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/object_models"
	"go-upcycle_connect-backend/app/models/package_models"
	"go-upcycle_connect-backend/app/models/payment_models"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"go-upcycle_connect-backend/utils/stripe"
	"math"
	"net/http"
	"time"
)

type createPaymentDTO struct {
	ObjectId   string `json:"object_id"`
	LockerId   string `json:"locker_id"`
	SuccessURL string `json:"success_url"`
	CancelURL  string `json:"cancel_url"`
}

const deliveryCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func genDeliveryCode() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = deliveryCodeAlphabet[int(b[i])%len(deliveryCodeAlphabet)]
	}
	return string(b)
}

// CreatePaymentHandler — POST /payments/checkout (auth required)
// Creates a one-time Stripe Checkout Session to buy the annonce (object) and
// returns its hosted URL. The amount is derived from the object price server
// side — the client never sends a price.
func CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	buyerId := auth_middleware.GetUserId(r.Context())

	var dto createPaymentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	if dto.ObjectId == "" {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}

	var object object_models.Object
	if err := object.Get([]string{"id", "name", "price", "user_id"}, db.IdClause, dto.ObjectId); err != nil {
		response.NewErrorMessage(w, response.ErrObjectNotFound, http.StatusNotFound)
		return
	}

	// A free annonce (don) or a zero/negative price has nothing to pay.
	if object.Price <= 0 {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	// The seller cannot buy their own annonce.
	if object.UserId == buyerId {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}

	amountCents := int64(math.Round(object.Price * 100))

	session, err := stripe.CreatePaymentCheckoutSession(
		object.Name,
		amountCents,
		dto.SuccessURL,
		dto.CancelURL,
		buyerId,
		map[string]string{
			"type":      "object_purchase",
			"object_id": object.Id,
			"locker_id": dto.LockerId,
		},
	)
	if err != nil {
		log.Info("stripe create payment session: " + err.Error())
		response.NewErrorMessage(w, response.ErrStripe, http.StatusBadGateway)
		return
	}

	// Record a pending payment so the webhook can reconcile it to "paid".
	_ = payment_models.Upsert(payment_models.Payment{
		ObjectId:        object.Id,
		UserId:          buyerId,
		StripeSessionId: session.ID,
		AmountCents:     int(amountCents),
		Status:          "pending",
	})

	// Livraison en casier : deux codes (depot vendeur + retrait acheteur).
	if dto.LockerId != "" {
		expiry := time.Now().Add(7 * 24 * time.Hour).Format("2006-01-02 15:04:05")
		package_models.CreateDelivery(object.Id, dto.LockerId, buyerId, genDeliveryCode(), genDeliveryCode(), session.ID, expiry)
	}

	response.NewSuccessData(w, map[string]string{"url": session.URL})
}

// GetPaymentStatusHandler — GET /payments/session/{id} (auth required)
// Reads the session from Stripe, reflects the result in DB, returns the status.
// Used by the success page to poll while the webhook confirms the payment.
func GetPaymentStatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusNotFound)
		return
	}

	session, err := stripe.GetCheckoutSession(id)
	if err != nil {
		log.Info("stripe get payment session: " + err.Error())
		response.NewErrorMessage(w, response.ErrStripe, http.StatusBadGateway)
		return
	}

	status := "pending"
	switch session.PaymentStatus {
	case "paid", "no_payment_required":
		status = "paid"
		_ = payment_models.Upsert(payment_models.Payment{
			ObjectId:        session.Metadata["object_id"],
			UserId:          session.ClientReferenceID,
			StripeSessionId: session.ID,
			AmountCents:     int(session.AmountTotal),
			Status:          "paid",
		})
		package_models.MarkPaidBySession(session.ID)
	case "unpaid":
		status = "unpaid"
	}

	response.NewSuccessData(w, map[string]string{
		"status":         status,
		"customer_email": session.CustomerDetails.Email,
	})
}

package payment_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/object_models"
	"go-upcycle_connect-backend/app/models/payment_models"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"go-upcycle_connect-backend/utils/stripe"
	"math"
	"net/http"
)

type createPaymentDTO struct {
	ObjectId   string `json:"object_id"`
	SuccessURL string `json:"success_url"`
	CancelURL  string `json:"cancel_url"`
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
	case "unpaid":
		status = "unpaid"
	}

	response.NewSuccessData(w, map[string]string{
		"status":         status,
		"customer_email": session.CustomerDetails.Email,
	})
}

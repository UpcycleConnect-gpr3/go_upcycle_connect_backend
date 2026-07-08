package billing_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/package_models"
	"go-upcycle_connect-backend/app/models/payment_models"
	"go-upcycle_connect-backend/app/models/subscription_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
	"go-upcycle_connect-backend/utils/stripe"
	"io"
	"net/http"
	"os"
)

type createSessionDTO struct {
	PriceID    string `json:"price_id"`
	Mode       string `json:"mode"`
	SuccessURL string `json:"success_url"`
	CancelURL  string `json:"cancel_url"`
}

func CreateCheckoutSessionHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	userId := auth_middleware.GetUserId(r.Context())

	var dto createSessionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	if dto.PriceID == "" {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}

	session, err := stripe.CreateCheckoutSession(dto.PriceID, dto.SuccessURL, dto.CancelURL, userId)
	if err != nil {
		log.Info("stripe create session: " + err.Error())
		response.NewErrorMessage(w, response.ErrStripe, http.StatusBadGateway)
		return
	}

	_ = subscription_models.Upsert(subscription_models.Subscription{
		UserId:          userId,
		StripeSessionId: session.ID,
		PriceId:         dto.PriceID,
		Status:          "pending",
	})

	response.NewSuccessData(w, map[string]string{"url": session.URL})
}

func GetCheckoutSessionHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusNotFound)
		return
	}

	session, err := stripe.GetCheckoutSession(id)
	if err != nil {
		log.Info("stripe get session: " + err.Error())
		response.NewErrorMessage(w, response.ErrStripe, http.StatusBadGateway)
		return
	}

	status := "pending"
	switch session.PaymentStatus {
	case "paid":
		status = "paid"
		_ = subscription_models.Upsert(subscription_models.Subscription{
			UserId:               session.ClientReferenceID,
			StripeSessionId:      session.ID,
			StripeSubscriptionId: session.Subscription,
			StripeCustomerId:     session.Customer,
			Status:               "active",
		})
	case "unpaid":
		status = "unpaid"
	}

	response.NewSuccessData(w, map[string]string{
		"status":         status,
		"customer_email": session.CustomerDetails.Email,
	})
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := stripe.VerifyWebhookSignature(
		payload, r.Header.Get("Stripe-Signature"), os.Getenv("STRIPE_WEBHOOK_SECRET"),
	); err != nil {
		log.Info("stripe webhook signature: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var event struct {
		Type string `json:"type"`
		Data struct {
			Object json.RawMessage `json:"object"`
		} `json:"data"`
	}
	if err := json.Unmarshal(payload, &event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "checkout.session.completed", "checkout.session.async_payment_succeeded":
		var s stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Object, &s); err == nil {
			fulfillCheckout(s)
		}
	case "checkout.session.expired", "checkout.session.async_payment_failed":
		var s stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Object, &s); err == nil {
			if s.Metadata["type"] == "object_purchase" {
				_ = payment_models.Upsert(payment_models.Payment{
					ObjectId:        s.Metadata["object_id"],
					UserId:          s.ClientReferenceID,
					StripeSessionId: s.ID,
					Status:          "failed",
				})
				package_models.CancelBySession(s.ID)
			}
		}
	case "customer.subscription.deleted":

	}

	w.WriteHeader(http.StatusOK)
}

func fulfillCheckout(s stripe.CheckoutSession) {
	if s.Metadata["type"] == "object_purchase" {
		_ = payment_models.Upsert(payment_models.Payment{
			ObjectId:        s.Metadata["object_id"],
			UserId:          s.ClientReferenceID,
			StripeSessionId: s.ID,
			AmountCents:     int(s.AmountTotal),
			Status:          "paid",
		})

		package_models.MarkPaidBySession(s.ID)
		return
	}

	_ = subscription_models.Upsert(subscription_models.Subscription{
		UserId:               s.ClientReferenceID,
		StripeSessionId:      s.ID,
		StripeSubscriptionId: s.Subscription,
		StripeCustomerId:     s.Customer,
		Status:               "active",
	})
}

// Package stripe is a tiny, dependency-free Stripe client built on the standard
// library. It only covers what the billing flow needs: create a Checkout
// Session, read one back, and verify webhook signatures. The Stripe secret key
// is read from the STRIPE_SECRET_KEY environment variable and never leaves here.
package stripe

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"encoding/json"
)

const apiBase = "https://api.stripe.com"

var httpClient = &http.Client{Timeout: 20 * time.Second}

func secretKey() string { return os.Getenv("STRIPE_SECRET_KEY") }

// CheckoutSession mirrors the subset of Stripe's Checkout Session we care about.
type CheckoutSession struct {
	ID                string `json:"id"`
	URL               string `json:"url"`
	PaymentStatus     string `json:"payment_status"`
	Customer          string `json:"customer"`
	Subscription      string `json:"subscription"`
	ClientReferenceID string `json:"client_reference_id"`
	CustomerDetails   struct {
		Email string `json:"email"`
	} `json:"customer_details"`
}

// CreateCheckoutSession creates a subscription Checkout Session and returns it
// (the hosted URL lives in session.URL). clientRef ties the session to a user.
func CreateCheckoutSession(priceID, successURL, cancelURL, clientRef string) (*CheckoutSession, error) {
	if secretKey() == "" {
		return nil, errors.New("STRIPE_SECRET_KEY is not set")
	}

	form := url.Values{}
	form.Set("mode", "subscription")
	form.Set("success_url", successURL)
	form.Set("cancel_url", cancelURL)
	form.Set("line_items[0][price]", priceID)
	form.Set("line_items[0][quantity]", "1")
	if clientRef != "" {
		form.Set("client_reference_id", clientRef)
	}

	req, err := http.NewRequest(http.MethodPost, apiBase+"/v1/checkout/sessions", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+secretKey())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return do(req)
}

// GetCheckoutSession fetches a Checkout Session by id.
func GetCheckoutSession(id string) (*CheckoutSession, error) {
	if secretKey() == "" {
		return nil, errors.New("STRIPE_SECRET_KEY is not set")
	}

	req, err := http.NewRequest(http.MethodGet, apiBase+"/v1/checkout/sessions/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+secretKey())

	return do(req)
}

func do(req *http.Request) (*CheckoutSession, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("stripe API error (%d): %s", resp.StatusCode, string(body))
	}

	var session CheckoutSession
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// VerifyWebhookSignature validates the Stripe-Signature header against the raw
// payload using the webhook signing secret (whsec_...). Same scheme as the
// official SDK: HMAC-SHA256 over "timestamp.payload".
func VerifyWebhookSignature(payload []byte, sigHeader, secret string) error {
	if secret == "" {
		return errors.New("missing STRIPE_WEBHOOK_SECRET")
	}

	var timestamp string
	var v1Signatures []string
	for _, part := range strings.Split(sigHeader, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			timestamp = kv[1]
		case "v1":
			v1Signatures = append(v1Signatures, kv[1])
		}
	}
	if timestamp == "" || len(v1Signatures) == 0 {
		return errors.New("invalid Stripe-Signature header")
	}

	// Reject events older than 5 minutes to mitigate replay attacks.
	if ts, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
		if time.Since(time.Unix(ts, 0)) > 5*time.Minute {
			return errors.New("timestamp outside tolerance")
		}
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "." + string(payload)))
	expected := hex.EncodeToString(mac.Sum(nil))

	for _, sig := range v1Signatures {
		if hmac.Equal([]byte(expected), []byte(sig)) {
			return nil
		}
	}
	return errors.New("signature mismatch")
}

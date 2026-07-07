package payment_models

import (
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"
	"math"

	"github.com/google/uuid"
)

const TABLE = "PAYMENTS"

// CommissionRate is the cut UpcycleConnect keeps on each annonce sale (10%).
const CommissionRate = 0.10

type Payment struct {
	Id              string `db:"id" json:"id"`
	ObjectId        string `db:"object_id" json:"object_id"`
	UserId          string `db:"user_id" json:"user_id"`
	StripeSessionId string `db:"stripe_session_id" json:"stripe_session_id"`
	AmountCents     int    `db:"amount_cents" json:"amount_cents"`
	CommissionCents int    `db:"commission_cents" json:"commission_cents"`
	NetCents        int    `db:"net_cents" json:"net_cents"`
	Currency        string `db:"currency" json:"currency"`
	Status          string `db:"status" json:"status"`
	CreatedAt       string `db:"created_at" json:"created_at"`
	UpdatedAt       string `db:"updated_at" json:"updated_at"`
}

// SplitCommission returns the platform commission and the seller's net share
// (in cents) for a given gross amount. commission = round(10%), net = rest.
func SplitCommission(amountCents int) (commission int, net int) {
	commission = int(math.Round(float64(amountCents) * CommissionRate))
	net = amountCents - commission
	return commission, net
}

func (m *Payment) Get(columns []string, by string, value ...any) error {
	return db.GetQuery[Payment](database.UpcycleConnect, TABLE, m, columns, by, value...)
}

// Upsert writes a payment keyed by stripe_session_id. It is idempotent so the
// same Stripe webhook can be delivered multiple times without creating
// duplicates or overwriting an already-known amount with a zero one.
func Upsert(p Payment) error {
	if p.Id == "" {
		p.Id = uuid.New().String()
	}
	if p.Currency == "" {
		p.Currency = "eur"
	}
	// Derive the 10% commission from the amount so it stays consistent wherever
	// the payment is written (creation, status poll, webhook).
	if p.AmountCents > 0 {
		p.CommissionCents, p.NetCents = SplitCommission(p.AmountCents)
	}
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (id, object_id, user_id, stripe_session_id, amount_cents, commission_cents, net_cents, currency, status, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW()) "+
			"ON DUPLICATE KEY UPDATE "+
			"object_id = COALESCE(NULLIF(VALUES(object_id), ''), object_id), "+
			"user_id = COALESCE(NULLIF(VALUES(user_id), ''), user_id), "+
			"amount_cents = GREATEST(VALUES(amount_cents), amount_cents), "+
			"commission_cents = GREATEST(VALUES(commission_cents), commission_cents), "+
			"net_cents = GREATEST(VALUES(net_cents), net_cents), "+
			"status = VALUES(status), "+
			"updated_at = NOW()",
		p.Id, p.ObjectId, p.UserId, p.StripeSessionId, p.AmountCents, p.CommissionCents, p.NetCents, p.Currency, p.Status,
	)
	if err != nil {
		log.Database("UPSERT "+TABLE, err)
	}
	return err
}

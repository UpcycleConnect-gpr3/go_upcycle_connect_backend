package subscription_models

import (
	"go-upcycle_connect-backend/var/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"

	"github.com/google/uuid"
)

const TABLE = "SUBSCRIPTIONS"

type Subscription struct {
	Id                   string `db:"id" json:"id"`
	UserId               string `db:"user_id" json:"user_id"`
	StripeSessionId      string `db:"stripe_session_id" json:"stripe_session_id"`
	StripeSubscriptionId string `db:"stripe_subscription_id" json:"stripe_subscription_id"`
	StripeCustomerId     string `db:"stripe_customer_id" json:"stripe_customer_id"`
	PriceId              string `db:"price_id" json:"price_id"`
	Status               string `db:"status" json:"status"`
	CreatedAt            string `db:"created_at" json:"created_at"`
	UpdatedAt            string `db:"updated_at" json:"updated_at"`
}

func (m *Subscription) Get(columns []string, by string, value ...any) error {
	return db.GetQuery[Subscription](database.UpcycleConnect, TABLE, m, columns, by, value...)
}

// GetByUser renvoie l'abonnement le plus recent de l'utilisateur (ou nil).
func GetByUser(userId string) *Subscription {
	var sub Subscription
	err := database.UpcycleConnect.Get(&sub,
		"SELECT id, user_id, stripe_session_id, stripe_subscription_id, stripe_customer_id, price_id, status, created_at, updated_at FROM "+TABLE+" WHERE user_id = ? ORDER BY created_at DESC LIMIT 1",
		userId,
	)
	if err != nil {
		return nil
	}
	return &sub
}

// Upsert writes a subscription keyed by stripe_session_id. It is idempotent so
// the same Stripe webhook can be delivered multiple times without harm.
func Upsert(s Subscription) error {
	if s.Id == "" {
		s.Id = uuid.New().String()
	}
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (id, user_id, stripe_session_id, stripe_subscription_id, stripe_customer_id, price_id, status, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW()) "+
			"ON DUPLICATE KEY UPDATE "+
			"stripe_subscription_id = VALUES(stripe_subscription_id), "+
			"stripe_customer_id = VALUES(stripe_customer_id), "+
			"price_id = COALESCE(NULLIF(VALUES(price_id), ''), price_id), "+
			"status = VALUES(status), "+
			"updated_at = NOW()",
		s.Id, s.UserId, s.StripeSessionId, s.StripeSubscriptionId, s.StripeCustomerId, s.PriceId, s.Status,
	)
	if err != nil {
		log.Database("UPSERT "+TABLE, err)
	}
	return err
}

package invoice_models

import (
	"os"

	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"
)

// Invoice est une vue "facture" derivee des paiements (achats d'annonces) et
// des abonnements de l'utilisateur. Le double conserve = l'enregistrement
// paiement/abonnement persiste ; le PDF est genere a la demande a partir de lui.
type Invoice struct {
	Ref         string `json:"ref"`  // "pay_<id>" ou "sub_<id>"
	Type        string `json:"type"` // "purchase" | "subscription"
	Label       string `json:"label"`
	AmountCents int    `json:"amount_cents"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
	IssuedAt    string `json:"issued_at"`
}

// subscriptionPlan renvoie (nom, montant en centimes) pour un price_id Stripe.
func subscriptionPlan(priceId string) (string, int) {
	switch priceId {
	case os.Getenv("STRIPE_PRICE_BUSINESS"):
		return "Business", 3000
	case os.Getenv("STRIPE_PRICE_BASIC"):
		return "Basic", 1500
	default:
		return "Pro", 3000
	}
}

type purchaseRow struct {
	Id          string `db:"id"`
	AmountCents int    `db:"amount_cents"`
	Currency    string `db:"currency"`
	Status      string `db:"status"`
	CreatedAt   string `db:"created_at"`
	ObjectName  string `db:"object_name"`
}

type subscriptionRow struct {
	Id        string `db:"id"`
	PriceId   string `db:"price_id"`
	Status    string `db:"status"`
	CreatedAt string `db:"created_at"`
}

// GetUserInvoices liste les factures de l'utilisateur (plus recentes d'abord).
func GetUserInvoices(userId string) []Invoice {
	invoices := []Invoice{}

	purchases := []purchaseRow{}
	err := database.UpcycleConnect.Select(&purchases,
		"SELECT p.id, p.amount_cents, COALESCE(p.currency,'eur') AS currency, p.status, p.created_at, "+
			"COALESCE(o.name,'Annonce') AS object_name "+
			"FROM PAYMENTS p LEFT JOIN OBJECTS o ON p.object_id = o.id "+
			"WHERE p.user_id = ? AND p.status = 'paid' ORDER BY p.created_at DESC",
		userId,
	)
	if err != nil {
		log.Database("SELECT USER PAYMENTS FOR INVOICES", err)
	}
	for _, p := range purchases {
		invoices = append(invoices, Invoice{
			Ref:         "pay_" + p.Id,
			Type:        "purchase",
			Label:       "Achat — " + p.ObjectName,
			AmountCents: p.AmountCents,
			Currency:    p.Currency,
			Status:      p.Status,
			IssuedAt:    p.CreatedAt,
		})
	}

	subs := []subscriptionRow{}
	err = database.UpcycleConnect.Select(&subs,
		"SELECT id, COALESCE(price_id,'') AS price_id, status, created_at FROM SUBSCRIPTIONS "+
			"WHERE user_id = ? ORDER BY created_at DESC",
		userId,
	)
	if err != nil {
		log.Database("SELECT USER SUBSCRIPTIONS FOR INVOICES", err)
	}
	for _, s := range subs {
		plan, amount := subscriptionPlan(s.PriceId)
		invoices = append(invoices, Invoice{
			Ref:         "sub_" + s.Id,
			Type:        "subscription",
			Label:       "Abonnement " + plan,
			AmountCents: amount,
			Currency:    "eur",
			Status:      s.Status,
			IssuedAt:    s.CreatedAt,
		})
	}

	return invoices
}

// GetUserInvoiceByRef renvoie une facture precise de l'utilisateur, ou nil si
// elle n'existe pas / ne lui appartient pas.
func GetUserInvoiceByRef(userId, ref string) *Invoice {
	for _, inv := range GetUserInvoices(userId) {
		if inv.Ref == ref {
			return &inv
		}
	}
	return nil
}

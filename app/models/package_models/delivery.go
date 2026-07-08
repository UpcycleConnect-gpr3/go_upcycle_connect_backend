package package_models

import (
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"

	"github.com/google/uuid"
)

// CreateDelivery cree une livraison en casier suite a un achat : deux codes
// (depot pour le vendeur, retrait pour l'acheteur), liee a la session Stripe.
// Statut initial 'awaiting_payment' (bascule en 'awaiting_deposit' au paiement).
func CreateDelivery(objectId, lockerId, buyerId, depositCode, retrieveCode, stripeSession, expiry string) *Package {
	id := uuid.New().String()
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (id, code, retrieve_code, locker_id, object_id, buyer_id, stripe_session_id, status, expiry_date, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, 'awaiting_payment', ?, NOW(), NOW())",
		id, depositCode, retrieveCode, lockerId, objectId, buyerId, stripeSession, expiry,
	)
	if err != nil {
		log.Database("INSERT DELIVERY PACKAGE", err)
		return nil
	}
	return &Package{Id: id}
}

// MarkPaidBySession bascule une livraison en 'awaiting_deposit' apres paiement
// (idempotent : ne touche que les livraisons encore en attente de paiement).
func MarkPaidBySession(sessionId string) {
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET status = 'awaiting_deposit', updated_at = NOW() "+
			"WHERE stripe_session_id = ? AND status = 'awaiting_payment'",
		sessionId,
	)
	if err != nil {
		log.Database("MARK DELIVERY PAID", err)
	}
}

// CancelBySession annule une livraison dont le paiement a echoue/expire.
func CancelBySession(sessionId string) {
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET status = 'canceled', updated_at = NOW() "+
			"WHERE stripe_session_id = ? AND status = 'awaiting_payment'",
		sessionId,
	)
	if err != nil {
		log.Database("CANCEL DELIVERY", err)
	}
}

// MarkDeposited passe la livraison en 'deposited' (le vendeur a depose l'objet).
func MarkDeposited(id string) error {
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET status = 'deposited', updated_at = NOW() WHERE id = ?",
		id,
	)
	if err != nil {
		log.Database("MARK DELIVERY DEPOSITED", err)
	}
	return err
}

// GetByAnyCode renvoie le package dont le code (depot) OU le code de retrait
// correspond. Sert au dépôt direct comme à la livraison d'achat.
func GetByAnyCode(code string) *Package {
	var pkg Package
	err := database.UpcycleConnect.Get(&pkg,
		"SELECT "+packageColumns+" FROM "+TABLE+" WHERE code = ? OR retrieve_code = ?",
		code, code,
	)
	if err != nil {
		return nil
	}
	return &pkg
}

// DeliverySummary : une livraison en casier, cote vendeur ou acheteur.
type DeliverySummary struct {
	PackageId    string  `db:"package_id" json:"package_id"`
	Code         string  `db:"code" json:"code"`
	RetrieveCode string  `db:"retrieve_code" json:"retrieve_code"`
	ObjectId     string  `db:"object_id" json:"object_id"`
	ObjectName   string  `db:"object_name" json:"object_name"`
	Price        float64 `db:"price" json:"price"`
	LockerName   string  `db:"locker_name" json:"locker_name"`
	LockerCity   string  `db:"locker_city" json:"locker_city"`
	Status       string  `db:"status" json:"status"`
	ExpiryDate   string  `db:"expiry_date" json:"expiry_date"`
}

const deliverySelect = "p.id AS package_id, p.code AS code, COALESCE(p.retrieve_code,'') AS retrieve_code, " +
	"o.id AS object_id, o.name AS object_name, o.price AS price, " +
	"l.name AS locker_name, l.city AS locker_city, p.status AS status, COALESCE(p.expiry_date,'') AS expiry_date " +
	"FROM " + TABLE + " p JOIN OBJECTS o ON p.object_id = o.id JOIN LOCKERS l ON p.locker_id = l.id "

// GetSellerDeliveries : livraisons a deposer par le vendeur (objets lui
// appartenant, statut 'awaiting_deposit') — expose le code de depot.
func GetSellerDeliveries(sellerId string) []DeliverySummary {
	result := []DeliverySummary{}
	err := database.UpcycleConnect.Select(&result,
		"SELECT "+deliverySelect+"WHERE o.user_id = ? AND p.status = 'awaiting_deposit' ORDER BY p.created_at DESC",
		sellerId,
	)
	if err != nil {
		log.Database("SELECT SELLER DELIVERIES", err)
	}
	return result
}

// GetBuyerDeliveries : achats a recuperer par l'acheteur (statut 'deposited')
// — expose le code de retrait.
func GetBuyerDeliveries(buyerId string) []DeliverySummary {
	result := []DeliverySummary{}
	err := database.UpcycleConnect.Select(&result,
		"SELECT "+deliverySelect+"WHERE p.buyer_id = ? AND p.status = 'deposited' ORDER BY p.created_at DESC",
		buyerId,
	)
	if err != nil {
		log.Database("SELECT BUYER DELIVERIES", err)
	}
	return result
}

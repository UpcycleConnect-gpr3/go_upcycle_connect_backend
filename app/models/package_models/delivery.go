package package_models

import (
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"

	"github.com/google/uuid"
)

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

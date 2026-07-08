package package_models

import (
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"

	"github.com/google/uuid"
)

const packageColumns = "id, COALESCE(weight,0) AS weight, code, locker_id, COALESCE(order_id, '') AS order_id, COALESCE(object_id, '') AS object_id, status, COALESCE(expiry_date, '') AS expiry_date, COALESCE(retrieve_code,'') AS retrieve_code, COALESCE(buyer_id,'') AS buyer_id, COALESCE(stripe_session_id,'') AS stripe_session_id, created_at, updated_at"

func CreateDeposit(objectId, lockerId, code string, weight int, expiryDate string) *Package {
	id := uuid.New().String()
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (id, weight, code, locker_id, object_id, status, expiry_date, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?, 'deposited', ?, NOW(), NOW())",
		id, weight, code, lockerId, objectId, expiryDate,
	)
	if err != nil {
		log.Database("INSERT DEPOSIT PACKAGE", err)
		return nil
	}
	return &Package{Id: id, Code: code, ObjectId: objectId, LockerId: lockerId, Status: "deposited", ExpiryDate: expiryDate}
}

func GetByCode(code string) *Package {
	var pkg Package
	err := database.UpcycleConnect.Get(&pkg,
		"SELECT "+packageColumns+" FROM "+TABLE+" WHERE code = ?",
		code,
	)
	if err != nil {
		return nil
	}
	return &pkg
}

func MarkRetrieved(id string) error {
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET status = 'retrieved', updated_at = NOW() WHERE id = ?",
		id,
	)
	if err != nil {
		log.Database("MARK PACKAGE RETRIEVED", err)
	}
	return err
}

type DepositedSummary struct {
	PackageId  string  `db:"package_id" json:"package_id"`
	Code       string  `db:"code" json:"code"`
	ObjectId   string  `db:"object_id" json:"object_id"`
	ObjectName string  `db:"object_name" json:"object_name"`
	Category   string  `db:"category" json:"category"`
	Score      float64 `db:"score" json:"score"`
	ImagePath  string  `db:"image_path" json:"image_path"`
	LockerName string  `db:"locker_name" json:"locker_name"`
	LockerCity string  `db:"locker_city" json:"locker_city"`
	ExpiryDate string  `db:"expiry_date" json:"expiry_date"`
}

func GetDepositedPackages() []DepositedSummary {
	result := []DepositedSummary{}
	err := database.UpcycleConnect.Select(&result,
		"SELECT p.id AS package_id, p.code AS code, o.id AS object_id, o.name AS object_name, "+
			"o.category AS category, o.score AS score, o.image_path AS image_path, "+
			"l.name AS locker_name, l.city AS locker_city, COALESCE(p.expiry_date,'') AS expiry_date "+
			"FROM "+TABLE+" p "+
			"JOIN OBJECTS o ON p.object_id = o.id "+
			"JOIN LOCKERS l ON p.locker_id = l.id "+
			"WHERE p.status = 'deposited' ORDER BY p.created_at DESC",
	)
	if err != nil {
		log.Database("SELECT DEPOSITED PACKAGES", err)
	}
	return result
}

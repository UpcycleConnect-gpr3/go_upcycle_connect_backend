package package_models

import (
	"go-upcycle_connect-backend/var/database"
	"go-upcycle_connect-backend/utils/log"

	"github.com/google/uuid"
)

const packageColumns = "id, weight, code, locker_id, COALESCE(order_id, '') AS order_id, COALESCE(object_id, '') AS object_id, status, COALESCE(expiry_date, '') AS expiry_date, created_at, updated_at"

// CreateDeposit cree un package depose (status 'deposited') pour un objet, avec
// un code de recuperation et une date d'expiration.
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

// GetByCode renvoie le package correspondant au code, ou nil.
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

// MarkRetrieved passe le package en 'retrieved'.
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

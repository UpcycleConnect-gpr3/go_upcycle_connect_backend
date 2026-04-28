package package_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"

	"github.com/google/uuid"
)

const TABLE = "PACKAGES"

type Package struct {
	Id        string `db:"id" json:"id"`
	Weight    int    `db:"weight" json:"weight"`
	Code      string `db:"code" json:"code"`
	LockerId  string `db:"locker_id" json:"locker_id"`
	OrderId   string `db:"order_id" json:"order_id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type CreatePackageDTO struct {
	Weight   int
	Code     string
	LockerId string
	OrderId  string
}

type UpdatePackageDTO struct {
	Weight   int
	Code     string
	LockerId string
	OrderId  string
}

func (m *Package) Get(columns []string, by string, value any) error {
	return db.GetQuery[Package](database.UpcycleConnect, TABLE, columns, by, value, m)
}

func (m *Package) All(columns []string, dest *[]Package) error {
	return db.AllQuery[Package](database.UpcycleConnect, TABLE, columns, dest)
}

func CreatePackage(dto CreatePackageDTO) *Package {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Code)
	packageId := uuid.New().String()
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (id, weight, code, locker_id, order_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())",
		packageId, dto.Weight, dto.Code, dto.LockerId, dto.OrderId,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Package{Id: packageId}
}

func UpdatePackage(id string, dto UpdatePackageDTO) *Package {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %s", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET weight=?, code=?, locker_id=?, order_id=?, updated_at=NOW() WHERE id=?",
		dto.Weight, dto.Code, dto.LockerId, dto.OrderId, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Package{Id: id}
}

func DeletePackage(id string) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %s", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

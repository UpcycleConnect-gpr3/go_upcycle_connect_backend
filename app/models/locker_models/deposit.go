package locker_models

import (
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"
)

const lockerColumns = "id, name, street, city, zip_code, capacity, available_slots, created_at, updated_at"

func GetAvailableLockers(city string) []Locker {
	lockers := []Locker{}
	var err error
	if city != "" {
		err = database.UpcycleConnect.Select(&lockers,
			"SELECT "+lockerColumns+" FROM "+TABLE+" WHERE available_slots > 0 AND city = ? ORDER BY city",
			city,
		)
	} else {
		err = database.UpcycleConnect.Select(&lockers,
			"SELECT "+lockerColumns+" FROM "+TABLE+" WHERE available_slots > 0 ORDER BY city",
		)
	}
	if err != nil {
		log.Database("SELECT AVAILABLE LOCKERS", err)
	}
	return lockers
}

func AdjustSlots(id string, delta int) error {
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET available_slots = LEAST(capacity, GREATEST(0, available_slots + ?)), updated_at = NOW() WHERE id = ?",
		delta, id,
	)
	if err != nil {
		log.Database("ADJUST LOCKER SLOTS", err)
	}
	return err
}

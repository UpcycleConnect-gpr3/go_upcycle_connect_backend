package config

import (
	"go-upcycle_connect-backend/internal"
	"go-upcycle_connect-backend/var/database"
	"os"
)

func InitDatabase() {

	database.UpcycleConnect = internal.NewDatabase(
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))
}

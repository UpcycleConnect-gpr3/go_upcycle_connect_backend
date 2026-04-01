package database

import (
	"go-upcycle_connect-backend/config"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/internal"
	"go-upcycle_connect-backend/utils/log"

	"github.com/joho/godotenv"
)

func initialize() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Config Initialization
	config.InitDatabase()

	err = database.UpcycleConnect.Ping()

	if err != nil {
		log.Fatal(err)
	}

	internal.CreateTableMigrations(database.UpcycleConnect)

}

func Migrate() {

	initialize()

	internal.Migrate(database.UpcycleConnect)

}

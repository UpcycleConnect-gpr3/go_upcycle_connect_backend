package database

import (
	"go-upcycle_connect-backend/config"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/migration"
	"go-upcycle_connect-backend/var/database"

	"github.com/joho/godotenv"
)

func initialize() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config.InitDatabase()

	err = database.UpcycleConnect.Ping()

	if err != nil {
		log.Fatal(err)
	}

	migration.CreateTableMigrations(database.UpcycleConnect)

}

func Migrate() {

	initialize()

	migration.Migrate(database.UpcycleConnect)

}

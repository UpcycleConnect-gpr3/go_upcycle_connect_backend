package internal

import (
	"fmt"
	"go-upcycle_connect-backend/utils/log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

const MigrationsPath = "database/migrations/"

func Migrate(db *sqlx.DB) {
	migrations, errorToGetNotExistingMigrations := getNotExistingMigrations(db)
	if errorToGetNotExistingMigrations != nil {
		log.Fatal(errorToGetNotExistingMigrations)
	}

	ShowStagingMigrations(migrations)

	if len(migrations) == 0 {
		return
	}

	var migrate string

	fmt.Println("")
	fmt.Println("-----------------------------")
	fmt.Print("Do you want to migrate? [Y/n] : ")
	_, errorToRead := fmt.Scanln(&migrate)

	if errorToRead != nil {
		log.Fatal(errorToRead)
		return
	}

	switch strings.ToLower(migrate) {
	case "y", "yes":
		executeMigration(db, migrations)
	case "n", "no":
		return
	}
}

func CreateTableMigrations(db *sqlx.DB) {
	systemTable, errorToLoadMigrationFile := readMigrationFile("system/migration_table.sql")

	if errorToLoadMigrationFile != nil {
		log.Fatal(errorToLoadMigrationFile)
	}

	_, errorToCreateMigrationTable := db.Exec(systemTable)

	if errorToCreateMigrationTable != nil {
		log.Fatal(errorToCreateMigrationTable)
	}

	log.Info("Migration Table Ready")
}

func ShowStagingMigrations(migrations []string) {

	if len(migrations) == 0 {
		log.Info("No migrations found")
	}

	if len(migrations) > 0 {
		log.Info("Migration(s) founds !")
	}

	for _, migration := range migrations {
		log.Info("\t" + migration)
	}
}

func readMigrationFile(fileName string) (string, error) {
	fileData, err := os.ReadFile(fmt.Sprintf("%s%s", MigrationsPath, fileName))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(fileData), err
}

func getMigrationFiles() ([]string, error) {
	files, err := os.ReadDir(MigrationsPath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames, nil
}

func getExistingMigrations(db *sqlx.DB) ([]string, error) {
	fileNames, errorToGetMigration := getMigrationFiles()

	if errorToGetMigration != nil {
		log.Fatal(errorToGetMigration)
		return fileNames, errorToGetMigration
	}

	rows, errorToQuery := db.Query("SELECT migration FROM MIGRATIONS ORDER BY id")

	if errorToQuery != nil {
		log.Fatal(errorToQuery)
		return fileNames, errorToQuery
	}

	defer rows.Close()

	var existingMigrations []string

	for rows.Next() {
		var migration string
		if errorToScanRow := rows.Scan(&migration); errorToScanRow != nil {
			log.Fatal(errorToScanRow)
			return fileNames, errorToScanRow
		}
		existingMigrations = append(existingMigrations, migration)
	}

	return existingMigrations, nil
}

func getMigrationBatch(db *sqlx.DB) ([]int, error) {

	var batches []int

	rows, errorToQuery := db.Query("SELECT batch FROM MIGRATIONS GROUP BY batch ORDER BY batch")

	if errorToQuery != nil {
		log.Fatal(errorToQuery)
		return batches, errorToQuery
	}

	defer rows.Close()

	for rows.Next() {
		var batch int
		if errorToScanRow := rows.Scan(&batch); errorToScanRow != nil {
			log.Fatal(errorToScanRow)
			return batches, errorToScanRow
		}
		batches = append(batches, batch)
	}
	return batches, nil
}

func getNotExistingMigrations(db *sqlx.DB) ([]string, error) {
	fileNames, errorToGetMigration := getMigrationFiles()

	if errorToGetMigration != nil {
		log.Fatal(errorToGetMigration)
		return fileNames, errorToGetMigration
	}

	existingMigrations, errorToGetExistingMigrations := getExistingMigrations(db)

	if errorToGetExistingMigrations != nil {
		log.Fatal(errorToGetExistingMigrations)
		return fileNames, errorToGetExistingMigrations
	}

	existingMigrationsMap := make(map[string]bool)
	for _, migration := range existingMigrations {
		existingMigrationsMap[migration] = true
	}

	var notExistingMigrations []string
	for _, fileName := range fileNames {
		if strings.HasSuffix(fileName, ".up.sql") && !existingMigrationsMap[fileName] {
			notExistingMigrations = append(notExistingMigrations, fileName)
		}
	}

	return notExistingMigrations, nil
}

func executeMigration(db *sqlx.DB, migrations []string) {
	batches, errorToGetBatch := getMigrationBatch(db)

	var batch int

	if len(batches) == 0 {
		batch = 1
	} else {
		batch = batches[len(batches)-1] + 1
	}

	if errorToGetBatch != nil {
		log.Fatal(errorToGetBatch)
		return
	}

	for _, migration := range migrations {
		migrationContent, errorToLoadFile := readMigrationFile(migration)

		if errorToLoadFile != nil {
			log.Fatal(errorToLoadFile)
			return
		}

		_, errorToExecuteMigration := db.Exec(migrationContent)
		if errorToExecuteMigration != nil {
			log.Fatal(errorToExecuteMigration)
			return
		}

		fmt.Println(migration)

		_, errorToExecuteMigration = db.Exec("INSERT INTO MIGRATIONS (migration, batch) VALUES (?, ?)", migration, batch)
		if errorToExecuteMigration != nil {
			log.Fatal(errorToExecuteMigration)
			return
		}

	}
}

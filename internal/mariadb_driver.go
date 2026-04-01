package internal

import (
	"database/sql"
	"fmt"
	"go-upcycle_connect-backend/utils/log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func NewDatabase(user string, password string, host string, port string, dbname string) *sql.DB {
	var intPort, errorToConvert = strconv.Atoi(port)

	if errorToConvert != nil {
		log.Fatal(errorToConvert)
	}

	var sqlInfo = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, intPort, dbname)

	conn, err := sql.Open("mysql", sqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	log.Info(fmt.Sprintf("(CONFIG) Database Drive Initialized %s", dbname))

	return conn
}

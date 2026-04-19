package db

import (
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/sql_builder"

	"github.com/jmoiron/sqlx"
)

func GetQuery[T any](db *sqlx.DB, table string, columns []string, by string, value any, dest *T) error {
	query := sql_builder.SelectQuery(table, columns, by)
	err := db.Get(dest, query, value)
	if err != nil {
		log.Database(query, err)
		return err
	}
	return nil
}

func AllQuery[T any](db *sqlx.DB, table string, columns []string, dest *[]T) error {
	query := sql_builder.SelectsQuery(table, columns)
	err := db.Select(dest, query)
	if err != nil {
		log.Database(query, err)
		return err
	}
	return nil
}

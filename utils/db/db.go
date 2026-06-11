package db

import (
	"fmt"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/sql_builder"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	IdClause = "id = ?"
)

func GetQuery[T any](db *sqlx.DB, table string, dest *T, columns []string, by string, value ...any) error {
	query := sql_builder.SelectQuery(table, columns, by)
	err := db.Get(dest, query, value...)
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

func CreateQuery(db *sqlx.DB, table string, columns []string, args ...any) error {
	query := sql_builder.InsertQuery(table, columns)
	_, err := db.Exec(query, args...)
	if err != nil {
		log.Database(query, err)
		return err
	}
	return nil
}

func UpdateQuery[T any](db *sqlx.DB, table string, model T, idTag string) error {
	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type()

	var query = "UPDATE " + table + " SET "
	var args []any

	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		dbTag := modelType.Field(i).Tag.Get("db")

		if dbTag == "" || dbTag == idTag {
			continue
		}

		if field.Kind() != reflect.Pointer || field.IsNil() {
			continue
		}

		query += fmt.Sprintf("%s = ?, ", dbTag)
		args = append(args, field.Elem().Interface())
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE %s = ?", idTag)

	var idFieldValue reflect.Value
	for i := 0; i < modelValue.NumField(); i++ {
		if modelType.Field(i).Tag.Get("db") == idTag {
			idFieldValue = modelValue.Field(i)
			break
		}
	}

	if !idFieldValue.IsValid() {
		return fmt.Errorf("no field with db tag '%s' found", idTag)
	}

	args = append(args, idFieldValue.Interface())

	_, err := db.Exec(query, args...)
	if err != nil {
		log.Database(query, err)
		return err
	}
	return nil
}

func DeleteQuery(db *sqlx.DB, table string, whereClause string, args ...any) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, whereClause)
	_, err := db.Exec(query, args...)
	if err != nil {
		log.Database(query, err)
		return err
	}
	return nil
}

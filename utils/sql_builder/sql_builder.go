package sql_builder

import (
	"fmt"
	"strings"
)

func SelectQuery(table string, columns []string, whereClause string) string {
	cols := "*"
	if len(columns) > 0 {
		cols = strings.Join(columns, ", ")
	}
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s", cols, table, whereClause)
}

func SelectsQuery(table string, columns []string) string {
	cols := "*"
	if len(columns) > 0 {
		cols = strings.Join(columns, ", ")
	}
	return fmt.Sprintf("SELECT %s FROM %s", cols, table)
}

func InsertQuery(table string, columns []string) string {
	placeholders := strings.Repeat("?, ", len(columns))
	placeholders = strings.TrimSuffix(placeholders, ", ")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), placeholders)
}

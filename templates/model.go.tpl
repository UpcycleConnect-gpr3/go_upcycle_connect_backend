package {{.PackageName}}

import (
	"{{.ModuleName}}/database"
	"{{.ModuleName}}/utils/db"
)

const TABLE = "{{.ResourceUpper}}S"

type {{.ResourceName}} struct {
	// TODO: Add fields
}

type Create{{.ResourceName}}DTO struct {
	// TODO: Add fields
}

type Update{{.ResourceName}}DTO struct {
	// TODO: Add fields
}

func (m *{{.ResourceName}}) Get(columns []string, by string, value any) error {
	return db.GetQuery[{{.ResourceName}}](database.Auth, TABLE, columns, by, value, m)
}

func (m *{{.ResourceName}}) All(columns []string, dest *[]{{.ResourceName}}) error {
	return db.AllQuery[{{.ResourceName}}](database.Auth, TABLE, columns, dest)
}

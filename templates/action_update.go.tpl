package {{.PackageName}}

import (
	"{{.ModuleName}}/app/models/{{.ResourceLower}}_models"
	"{{.ModuleName}}/utils/rules"
)

type Update{{.ResourceName}}DTO struct {
	// TODO: Add fields
}

func Update{{.ResourceName}}(id int, dto Update{{.ResourceName}}DTO) ([]rules.ValidationError, *{{.ResourceLower}}_models.{{.ResourceName}}) {
	var errs []rules.ValidationError

	// TODO: Add validation rules
	// rules.StringMinLength(dto.Field, 1, "field", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	{{.ResourceLower}} := {{.ResourceLower}}_models.Update{{.ResourceName}}(id, dto)

	return nil, {{.ResourceLower}}
}

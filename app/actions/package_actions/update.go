package package_actions

import (
	"go-upcycle_connect-backend/app/models/package_models"
	"go-upcycle_connect-backend/utils/rules"
)

func UpdatePackage(id string, dto package_models.UpdatePackageDTO) ([]rules.ValidationError, *package_models.Package) {
	var errs []rules.ValidationError

	rules.IntMinLength(dto.Weight, 0, "weight", &errs)
	rules.StringMinLength(dto.Code, 1, "code", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	pkg := package_models.UpdatePackage(id, dto)

	return nil, pkg
}

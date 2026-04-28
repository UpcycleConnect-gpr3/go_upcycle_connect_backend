package locker_actions

import (
	"go-upcycle_connect-backend/app/models/locker_models"
	"go-upcycle_connect-backend/utils/rules"
)

func UpdateLocker(id string, dto locker_models.UpdateLockerDTO) ([]rules.ValidationError, *locker_models.Locker) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMinLength(dto.Street, 1, "street", &errs)
	rules.StringMinLength(dto.City, 1, "city", &errs)
	rules.StringMinLength(dto.ZipCode, 1, "zip_code", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	locker := locker_models.UpdateLocker(id, dto)

	return nil, locker
}

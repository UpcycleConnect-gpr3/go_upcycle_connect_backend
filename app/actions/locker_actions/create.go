package locker_actions

import (
	"go-upcycle_connect-backend/app/models/locker_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateLocker(dto locker_models.CreateLockerDTO) ([]rules.ValidationError, *locker_models.Locker) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMinLength(dto.Street, 1, "street", &errs)
	rules.StringMinLength(dto.City, 1, "city", &errs)
	rules.StringMinLength(dto.ZipCode, 1, "zip_code", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	locker := locker_models.CreateLocker(dto)

	return nil, locker
}

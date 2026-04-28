package step_actions

import (
	"go-upcycle_connect-backend/app/models/step_models"
	"go-upcycle_connect-backend/utils/rules"
)

func UpdateStep(id string, dto step_models.UpdateStepDTO) ([]rules.ValidationError, *step_models.Step) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	step := step_models.UpdateStep(id, dto)

	return nil, step
}

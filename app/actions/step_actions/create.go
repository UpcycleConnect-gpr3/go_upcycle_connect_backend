package step_actions

import (
	"go-upcycle_connect-backend/app/models/step_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateStep(dto step_models.CreateStepDTO) ([]rules.ValidationError, *step_models.Step) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	step := step_models.CreateStep(dto)

	return nil, step
}

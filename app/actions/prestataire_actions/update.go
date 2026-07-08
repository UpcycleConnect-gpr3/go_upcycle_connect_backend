package prestataire_actions

import (
	"go-upcycle_connect-backend/app/models/prestataire_models"
	"go-upcycle_connect-backend/utils/rules"
)

func UpdatePrestataire(id int, dto prestataire_models.UpdatePrestataireDTO) ([]rules.ValidationError, *prestataire_models.Prestataire) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	prestataire := prestataire_models.UpdatePrestataire(id, dto)
	return nil, prestataire
}

package ad_actions

import (
	"go-upcycle_connect-backend/app/models/ad_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateAd(dto ad_models.CreateAdDTO) ([]rules.ValidationError, *ad_models.Ad) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)

	if dto.Status == "" {
		dto.Status = "active"
	}

	if len(errs) > 0 {
		return errs, nil
	}

	ad := ad_models.CreateAd(dto)
	return nil, ad
}

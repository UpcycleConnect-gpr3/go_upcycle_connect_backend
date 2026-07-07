package user_actions

import (
	"go-upcycle_connect-backend/app/models/user_models"
	"go-upcycle_connect-backend/utils/rules"
)

func CreateUserFromToken(userId string) ([]rules.ValidationError, *user_models.User) {
	var errs []rules.ValidationError

	rules.StringMinLength(userId, 1, "userId", &errs)

	if len(errs) > 0 {
		return errs, nil
	}

	// username et email sont NOT NULL UNIQUE : l'id du token (unique par
	// construction) sert de valeur de provisionnement, le profil pourra
	// etre complete ensuite via PATCH /auth/update.
	var user user_models.User
	err := user.Create(user_models.CreateUserDTO{
		Id:        userId,
		Username:  userId,
		Firstname: "",
		Lastname:  "",
		Email:     userId,
	})
	if err != nil {
		return nil, nil
	}

	return nil, &user
}

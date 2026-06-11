package auth_handlers

import (
	"go-upcycle_connect-backend/app/actions/user_actions"
	"go-upcycle_connect-backend/app/models/user_models"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/jwt"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		response.NewErrorMessage(w, response.ErrAuthTokenRequired, http.StatusUnauthorized)
		return
	}

	userId, err := jwt.VerifyJWT(tokenString)

	if err != nil {
		response.NewErrorMessage(w, response.ErrInvalidAuthToken, http.StatusUnauthorized)
		return
	}

	var existingUser user_models.User
	err = existingUser.Get([]string{"id"}, db.IdClause, userId)
	if err == nil {

		columns := []string{"id", "firstname", "lastname", "email", "created_at", "updated_at"}
		existingUser.Get(columns, db.IdClause, userId)
		response.NewSuccessData(w, existingUser)
		return
	}

	validationErrors, newUser := user_actions.CreateUserFromToken(userId)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if newUser == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}

	var createdUser user_models.User
	columns := []string{"id", "firstname", "lastname", "email", "created_at", "updated_at"}
	createdUser.Get(columns, db.IdClause, userId)

	response.NewSuccessData(w, createdUser)
}

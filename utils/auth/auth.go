package auth

import (
	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/user_models"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

type AuthId struct {
	UserId string
}

func Auth(r *http.Request) *AuthId {
	return &AuthId{UserId: auth_middleware.GetUserId(r.Context())}
}

func (authId *AuthId) User(w http.ResponseWriter, columns []string) *user_models.User {
	columns = append([]string{"id"}, columns...)
	var user user_models.User
	err := user.Get(columns, "id = ?", authId.UserId)
	if err != nil {
		response.NewErrorMessage(w, response.ErrUserNotFound, http.StatusNotFound)
		return &user
	}
	return &user
}

func (authId *AuthId) Id() string {
	return authId.UserId
}

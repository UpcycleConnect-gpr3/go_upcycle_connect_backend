package auth_handlers

import (
	"encoding/json"
	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/user_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	userId := auth_middleware.GetUserId(r.Context())
	if userId == "" {
		response.NewErrorMessage(w, response.ErrAuthTokenRequired, http.StatusUnauthorized)
		return
	}

	var userDto user_models.UpdateUserDTO
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}

	var user user_models.User
	err = user.Update(userDto, userId)
	if err != nil {
		response.NewErrorMessage(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	response.NewSuccessData(w, user)
}

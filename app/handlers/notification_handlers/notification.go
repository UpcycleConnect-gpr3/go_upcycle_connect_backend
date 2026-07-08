package notification_handlers

import (
	"encoding/json"
	"net/http"

	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/activity_models"
	"go-upcycle_connect-backend/app/models/notification_models"
	"go-upcycle_connect-backend/utils/jwt"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
)

func MyNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())
	response.NewSuccessData(w, notification_models.GetUserNotifications(userId))
}

type createNotificationDTO struct {
	UserId string `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func StoreNotificationHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	if jwt.RoleFromToken(r.Header.Get("Authorization")) != "administrator" {
		response.NewErrorMessage(w, response.ErrForbidden, http.StatusForbidden)
		return
	}
	var dto createNotificationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}
	if dto.UserId == "" || dto.Title == "" {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}
	notif := notification_models.Create(dto.UserId, dto.Title, dto.Body)
	if notif == nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	activity_models.Record(auth_middleware.GetUserId(r.Context()), "notify", "user", dto.UserId, dto.Title)
	response.NewSuccessData(w, map[string]int{"id": notif.Id})
}

func MarkNotificationReadHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())
	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}
	if err := notification_models.MarkRead(id, userId); err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessMessage(w, "Notification marked as read")
}

package order_handlers

import (
	"net/http"

	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/order_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
)

func MyOrdersHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())
	response.NewSuccessData(w, order_models.GetUserOrders(userId))
}

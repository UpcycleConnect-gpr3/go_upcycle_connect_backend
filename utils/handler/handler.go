package handler

import (
	"go-upcycle_connect-backend/utils/response"
	"net/http"
	"strconv"
)

func ParsePathInt(w http.ResponseWriter, r *http.Request, key string, notFoundErr string) (int, bool) {
	id, err := strconv.Atoi(r.PathValue(key))
	if err != nil {
		response.NewErrorMessage(w, notFoundErr, http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

package request

import (
	"go-upcycle_connect-backend/utils/response"
	"net/http"
	"strconv"
)

type RequestValue struct {
	value string
}

func Request(r *http.Request, attribute string) *RequestValue {
	return &RequestValue{value: r.PathValue(attribute)}
}

func (requestValue *RequestValue) Value() string {
	return requestValue.value
}

func (requestValue *RequestValue) ConvertToInt(w http.ResponseWriter) int {
	value, err := strconv.Atoi(requestValue.value)
	if err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusBadRequest)
		return -1
	}
	return value
}

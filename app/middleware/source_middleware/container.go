package source_middleware

import (
	"go-upcycle_connect-backend/utils/response"
	"net/http"
)

func Container(allowedContainer string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			clientContainer := r.Header.Get("X-Container-Name")
			if clientContainer == "" {
				response.NewErrorMessage(w, "", http.StatusForbidden)
				return
			}
			if clientContainer != allowedContainer {
				response.NewErrorMessage(w, "", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}

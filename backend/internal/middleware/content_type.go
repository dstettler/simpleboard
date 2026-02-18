package middleware

import (
	"net/http"
	"strings"

	"cise.ufl.edu/no-frills-chess/pkg/response"
)

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			ct := r.Header.Get("Content-Type")
			if !strings.HasPrefix(ct, "application/json") {
				response.Error(w, http.StatusBadRequest, response.CodeInvalidInput, "Content-Type must be application/json")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

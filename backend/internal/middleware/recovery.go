package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"cise.ufl.edu/no-frills-chess/pkg/response"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().Interface("panic", err).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Msg("panic recovered")
				response.Error(w, http.StatusInternalServerError, response.CodeInternal, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

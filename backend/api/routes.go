package api

import (
	"github.com/go-chi/chi/v5"

	"cise.ufl.edu/no-frills-chess/internal/handler"
	"cise.ufl.edu/no-frills-chess/internal/middleware"
)

func NewRouter(corsOrigins []string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.CORS(corsOrigins))
	r.Use(middleware.Recovery)
	r.Use(middleware.Logger)
	r.Use(middleware.ContentTypeJSON)

	r.Get("/api/health", handler.Health)

	return r
}

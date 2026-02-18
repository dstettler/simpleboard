package api

import (
	"github.com/go-chi/chi/v5"

	"cise.ufl.edu/simpleboard/internal/handler"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/health", handler.Health)

	return r
}

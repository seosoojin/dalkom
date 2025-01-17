package handlers

import "github.com/go-chi/chi/v5"

type Http interface {
	RegisterRoutes(r *chi.Mux)
}

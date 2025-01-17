package cards

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/seosoojin/dalkom/internal/domain/handlers"
	"github.com/seosoojin/dalkom/internal/domain/pagination"
)

type Handler interface {
	GetCards(w http.ResponseWriter, r *http.Request)
	handlers.Http
}

type handler struct {
	service Service
}

var _ Handler = &handler{}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(r *chi.Mux) {
	r.Get("/cards", h.GetCards)
}

func (h *handler) GetCards(w http.ResponseWriter, r *http.Request) {
	page, err := pagination.NewPageFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cards, err := h.service.GetCards(r.Context(), page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, cards)
}

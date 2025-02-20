package cards

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/seosoojin/dalkom/internal/domain/handlers"
	"github.com/seosoojin/dalkom/internal/domain/pagination"
	"github.com/seosoojin/dalkom/pkg/models"
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
	r.Get("/cards/{id}", h.GetByID)
	r.Post("/cards", h.CreateCard)
}

func (h *handler) GetCards(w http.ResponseWriter, r *http.Request) {
	page, err := pagination.NewPageFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter, err := h.parseFilter(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cards, err := h.service.GetCards(r.Context(), filter, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, cards)
}

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	card, err := h.service.GetEnrichedCard(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, card)
}

func (h *handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	card := new(models.Card)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(b, card)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCard(r.Context(), card); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, card)
}

func (h *handler) parseFilter(r *http.Request) (map[string][]any, error) {
	filter := map[string][]any{}

	values := r.URL.Query()
	for key, value := range values {
		if key == "offset" || key == "limit" {
			continue
		}
		filter[key] = []any{}

		for _, v := range value {
			filter[key] = append(filter[key], v)
		}
	}

	return filter, nil
}

package binders

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/seosoojin/dalkom/internal/domain/auth"
	"github.com/seosoojin/dalkom/internal/domain/handlers"
	"github.com/seosoojin/dalkom/internal/domain/pagination"
	"github.com/seosoojin/dalkom/internal/gateways/middlewares"
	"github.com/seosoojin/dalkom/pkg/models"
)

type Handler interface {
	GetByUserID(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	handlers.Http
}

type handler struct {
	service        Service
	authMiddleware middlewares.Authenticator
}

var _ Handler = &handler{}

func NewHandler(service Service, auth middlewares.Authenticator) *handler {
	return &handler{
		service:        service,
		authMiddleware: auth,
	}
}

func (h *handler) RegisterRoutes(r *chi.Mux) {
	// r.Get("/users/{id}/binders", h.GetByUserID)
	r.Group(func(r chi.Router) {
		r.Use(h.authMiddleware.Authenticate())
		r.Get("/me/binders", h.GetByUserID)
		r.Get("/binders/{id}/cards", h.GetBinderCards)
		r.Post("/binders", h.Create)
		r.Put("/binders/{id}", h.Update)
		r.Patch("/binders/{id}/cards/{card_id}", h.AddCard)
		r.Delete("/binders/{id}/cards/{card_id}", h.RemoveCard)
		r.Get("/binders/{id}", h.GetByID)
		r.Delete("/binders/{id}", h.Delete)
	})

}

func (h *handler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	page, err := pagination.NewPageFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := auth.UserFromContext(r.Context())

	filter := h.parseFilter(r)

	binders, err := h.service.GetByUserID(r.Context(), user.ID, filter, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, binders)
}

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	binder, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := auth.UserFromContext(r.Context())

	if binder.UserID != user.ID {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	render.JSON(w, r, binder)
}

func (h *handler) GetBinderCards(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	binder, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := auth.UserFromContext(r.Context())

	if binder.UserID != user.ID {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	cards, err := h.service.GetBinderCards(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, cards)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	binder := new(models.Binder)

	err := json.NewDecoder(r.Body).Decode(binder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := auth.UserFromContext(r.Context())
	binder.UserID = user.ID

	if err := h.service.Create(r.Context(), binder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, binder)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	binder := new(models.Binder)

	err := json.NewDecoder(r.Body).Decode(binder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	binder.ID = id
	if err := h.service.Update(r.Context(), binder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, binder)
}

func (h *handler) AddCard(w http.ResponseWriter, r *http.Request) {
	binderID := chi.URLParam(r, "id")
	cardID := chi.URLParam(r, "card_id")

	if err := h.service.AddCard(r.Context(), binderID, cardID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, nil)
}

func (h *handler) RemoveCard(w http.ResponseWriter, r *http.Request) {
	binderID := chi.URLParam(r, "id")
	cardID := chi.URLParam(r, "card_id")

	if err := h.service.RemoveCard(r.Context(), binderID, cardID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, nil)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	binder, err := h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, binder)
}

func (h *handler) parseFilter(r *http.Request) map[string][]any {
	filter := map[string][]any{}

	values := r.URL.Query()
	for key, value := range values {
		if key == "page" || key == "limit" {
			continue
		}

		if key == "is_favorite" {
			isFavorite, err := strconv.ParseBool(value[0])
			if err != nil {
				continue
			}
			filter[key] = []any{isFavorite}
			continue
		}

		filter[key] = []any{}

		for _, v := range value {
			filter[key] = append(filter[key], v)
		}
	}

	return filter
}

package binders

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/seosoojin/dalkom/internal/domain/handlers"
	"github.com/seosoojin/dalkom/internal/domain/pagination"
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
	service Service
}

var _ Handler = &handler{}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(r *chi.Mux) {
	r.Get("/users/{id}/binders", h.GetByUserID)
	r.Get("/binders/{id}", h.GetByID)
	r.Post("/binders", h.Create)
	r.Put("/binders/{id}", h.Update)
	r.Delete("/binders/{id}", h.Delete)
}

func (h *handler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	page, err := pagination.NewPageFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := chi.URLParam(r, "user_id")
	binders, err := h.service.GetByUserID(r.Context(), userID, page)
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

	render.JSON(w, r, binder)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	binder := new(models.Binder)

	err := json.NewDecoder(r.Body).Decode(binder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	binder, err := h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, binder)
}

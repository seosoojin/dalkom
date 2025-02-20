package collections

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/seosoojin/dalkom/internal/domain/handlers"
	"github.com/seosoojin/dalkom/pkg/models"
)

type Handler interface {
	GetCollections(w http.ResponseWriter, r *http.Request)
	CreateCollection(w http.ResponseWriter, r *http.Request)
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
	r.Get("/collections", h.GetCollections)
	r.Post("/collections", h.CreateCollection)
}

func (h *handler) GetCollections(w http.ResponseWriter, r *http.Request) {
	Collections, err := h.service.GetCollections(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, Collections)
}

func (h *handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	collection := new(models.Collection)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(b, collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), collection); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, collection)
}

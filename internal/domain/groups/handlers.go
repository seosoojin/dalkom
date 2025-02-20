package groups

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
	GetGroups(w http.ResponseWriter, r *http.Request)
	CreateGroup(w http.ResponseWriter, r *http.Request)
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
	r.Get("/groups", h.GetGroups)
	r.Post("/groups", h.CreateGroup)
}

func (h *handler) GetGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.service.GetGroups(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, groups)
}

func (h *handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	group := new(models.Group)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(b, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), group); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, group)
}

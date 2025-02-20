package idols

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
	GetIdols(w http.ResponseWriter, r *http.Request)
	CreateIdol(w http.ResponseWriter, r *http.Request)
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
	r.Get("/idols", h.GetIdols)
	r.Post("/idols", h.CreateIdol)
}

func (h *handler) GetIdols(w http.ResponseWriter, r *http.Request) {
	idols, err := h.service.GetIdols(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, idols)
}

func (h *handler) CreateIdol(w http.ResponseWriter, r *http.Request) {
	idol := new(models.Idol)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(b, idol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), idol); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, idol)
}

package web

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/seosoojin/dalkom/internal/domain/handlers"
)

type Server struct {
	handlers []handlers.Http
	port     string
	router   *chi.Mux
}

func NewServer(port string, handlers ...handlers.Http) *Server {
	s := &Server{
		handlers: handlers,
		port:     port,
		router:   chi.NewRouter(),
	}

	s.router.Use(
		middleware.RequestID,
		middleware.Recoverer,
		middleware.RealIP,
		middleware.Logger,
		middleware.URLFormat,
		middleware.GetHead,
		middleware.Heartbeat("/health"),
	)

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	for _, h := range s.handlers {
		h.RegisterRoutes(s.router)
	}
}

func (s *Server) Run() {
	srv := http.Server{
		Addr:        ":" + s.port,
		Handler:     s.router,
		IdleTimeout: 10 * time.Second,
	}

	log.Println("Running on port", s.port)
	srv.ListenAndServe()
}

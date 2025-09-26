package api

import (
	"log"
	"net/http"
	"warehouse/api/handlers"
	"warehouse/api/middleware"
	"warehouse/internal/cfg"
	"warehouse/internal/services"
)

type Server struct {
	http *http.Server
}

func NewServer(cfg *cfg.Config, services *services.Services) *Server {
	mux := http.NewServeMux()
	registerHandlers(mux, services)
	handler := middleware.LogMiddleware(mux)

	s := &Server{&http.Server{
		Addr:    cfg.Server.Port,
		Handler: handler,
	}}

	return s
}

func (s *Server) Run() error {
	log.Printf("Server running on port %s \n", s.http.Addr)
	return s.http.ListenAndServe()
}

func registerHandlers(mux *http.ServeMux, services *services.Services) {
	h := handlers.NewHandler(services)

	mux.HandleFunc("/", h.GetDefault)
	mux.HandleFunc("/items", h.ItemHandler.HandleItems)
	mux.HandleFunc("/items/", h.ItemHandler.HandleItems)
}

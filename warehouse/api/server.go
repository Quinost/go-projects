package api

import (
	"log"
	"net/http"
	"warehouse/api/handlers"
	"warehouse/api/middleware"
	"warehouse/internal/repositories"
)

type Server struct {
	http *http.Server
}

func NewServer(repo *repositories.Repositories) *Server {
	mux := http.NewServeMux()
	registerHandlers(mux, repo)
	handler := middleware.LogMiddleware(mux)

	s := &Server{&http.Server{
		Addr:    ":8080",
		Handler: handler,
	}}

	return s
}

func (s *Server) Run() error {
	log.Printf("Server running on port %s \n", s.http.Addr)
	return s.http.ListenAndServe()
}

func registerHandlers(mux *http.ServeMux, repo *repositories.Repositories) {
	h := handlers.New(repo)

	mux.HandleFunc("/", h.GetDefault)
	mux.HandleFunc("/items", h.ItemHandler.HandleItems)
	mux.HandleFunc("/items/", h.ItemHandler.HandleItems)
}

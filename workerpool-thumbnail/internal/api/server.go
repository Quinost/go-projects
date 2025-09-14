package api

import (
	"context"
	"fmt"
	"net/http"
	"thumbnail/internal/services/thumbnail"
)

type Server struct {
	httpServer       *http.Server
	thumbnailService *thumbnail.ThumbnailService
}

func NewServer(thumbnailService *thumbnail.ThumbnailService) *Server {
	mux := http.NewServeMux()

	s := &Server{httpServer: &http.Server{
		Addr:    ":8080",
		Handler: mux,
	},
		thumbnailService: thumbnailService,
	}

	mux.HandleFunc("/upload", s.uploadHandler)
	mux.HandleFunc("/", s.defaultHandler)

	return s
}

func (s *Server) Start() error {
	fmt.Println("Started server on", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

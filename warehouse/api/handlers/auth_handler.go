package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"warehouse/internal/services"

	"github.com/google/uuid"
)

type AuthHandler struct {
	Handler
	service *services.JWTService
}

func NewAuthHandler(service *services.JWTService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/auth")
	method := r.Method

	switch {
	case strings.HasPrefix(path, "/login") && method == http.MethodPost:
		h.Login(w, r)
	default:
		writeError(w, http.StatusNotFound, "Not Found")
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	jwt, _ := h.service.GenerateJWT(uuid.New())

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "JWT %s", jwt)
}

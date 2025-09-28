package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"warehouse/internal/models"
	"warehouse/internal/services"
)

type AuthHandler struct {
	Handler
	jwtService  *services.JWTService
	userService *services.UserService
}

func NewAuthHandler(services *services.Services) *AuthHandler {
	return &AuthHandler{
		jwtService:  services.JWTService,
		userService: services.UserService,
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
	var loginReq models.LoginDto
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	user, err := h.userService.CheckAndGetUser(loginReq.Username, loginReq.Password)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	jwt, _ := h.jwtService.GenerateJWT(user.Id, user.Username)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", jwt)
}

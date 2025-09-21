package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"warehouse/internal/models"
	"warehouse/internal/repositories"

	"github.com/google/uuid"
)

type Handler struct {
	ItemHandler *ItemHandler
}

func New(repo *repositories.Repositories) *Handler {
	return &Handler{
		ItemHandler: &ItemHandler{repo: repo.ItemRep},
	}
}

func(h *Handler) extractID(path string) (uuid.UUID) {
    parts := strings.Split(path, "/")
    if len(parts) < 3 {
        return uuid.Nil
    }

	uid, err := uuid.Parse(parts[2])
	if err != nil {
		return uuid.Nil
	}
    return uid
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)

    resp := &models.Response{
        Error:  message,
        Status: models.ErrorStatus,
    }

    json.NewEncoder(w).Encode(resp)
}
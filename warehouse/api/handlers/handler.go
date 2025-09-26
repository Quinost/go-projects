package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"warehouse/internal/models"
	"warehouse/internal/services"

	"github.com/google/uuid"
)

type Handler struct {
	ItemHandler *ItemHandler
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		ItemHandler: NewItemHandler(services.ItemService),
	}
}

func (h *Handler) extractUUIDs(path string) []uuid.UUID {
	var ids []uuid.UUID
	segments := strings.SplitSeq(strings.Trim(path, "/"), "/")

	for segment := range segments {
		if id, err := uuid.Parse(segment); err == nil {
			ids = append(ids, id)
		}
	}

	return ids
}

func parseQueryParam[T any](query url.Values, key string) T {
	var result any
	var err error
	value := query.Get(key)
	if value == "" {
		return *new(T)
	}

	switch any(*new(T)).(type) {
	case int:
		result, err = strconv.Atoi(value)
	case float64:
		result, err = strconv.ParseFloat(value, 64)
	case bool:
		result, err = strconv.ParseBool(value)
	default:
		result = value
	}

	if err != nil {
		return *new(T)
	}

	return result.(T)
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

func writeOkJson(w http.ResponseWriter, object any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(object)
}

func writeCreatedJson(w http.ResponseWriter, object any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(object)
}

func writeOk(w http.ResponseWriter, status models.Status) {
	resp := &models.Response{
		Status: status,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

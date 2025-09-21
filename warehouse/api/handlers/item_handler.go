package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"warehouse/internal/models"
	"warehouse/internal/repositories"

	"github.com/google/uuid"
)

type ItemHandler struct {
	Handler
	repo *repositories.ItemRepository
}

func (h *ItemHandler) HandleItems(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/items")
	method := r.Method

	switch {
	case path == "" && method == http.MethodGet:
		h.GetAll(w, r)
	case strings.HasPrefix(path, "/") && method == http.MethodGet:
		h.GetById(w, r)
	case method == http.MethodPost:
		h.Add(w, r)
	case method == http.MethodPut:
		h.Update(w, r)
	case method == http.MethodDelete:
		h.Delete(w, r)
	default:
		writeError(w, http.StatusNotFound, "Not Found")
	}
}

func (h *ItemHandler) GetById(w http.ResponseWriter, r *http.Request) {
	uid := h.extractID(r.URL.Path)
	item, err := h.repo.GetById(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var page, limit int
	filter := query.Get("filter")

	if parsed, err := strconv.Atoi(query.Get("page")); err == nil {
		page = parsed - 1
	}

	if parsed, err := strconv.Atoi(query.Get("limit")); err == nil {
		limit = parsed
	}

	items, err := h.repo.GetAll(filter, page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func (h *ItemHandler) Add(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	item.Id = uuid.New()

	if err := h.repo.Add(&item); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to insert item")
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/items/%s", item.Id))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.NewResponse(item.Id.String(), models.CreatedStatus))
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	item.Id = h.extractID(r.URL.Path)

	if item.Id == uuid.Nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	existedItem, err := h.repo.GetById(item.Id)
	if err != nil || existedItem == nil {
		writeError(w, http.StatusBadRequest, "Item not found")
		return
	}

	if err := h.repo.Update(&item); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update item")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.NewResponse("", models.UpdatedStatus))
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid := h.extractID(r.URL.Path)

	if err := h.repo.Delete(uid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

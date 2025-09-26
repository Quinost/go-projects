package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"warehouse/internal/models"
	"warehouse/internal/services"
)

type ItemHandler struct {
	Handler
	service *services.ItemService
}

func NewItemHandler(service *services.ItemService) *ItemHandler {
	return &ItemHandler{
		service: service,
	}
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
	uid := h.extractUUIDs(r.URL.Path)[0]

	item, err := h.service.GetById(uid)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOkJson(w, item)
}

func (h *ItemHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filter := parseQueryParam[string](query, "filter")
	page := parseQueryParam[int](query, "page") - 1
	limit := parseQueryParam[int](query, "limit")

	items, err := h.service.GetAll(filter, page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeOkJson(w, items)
}

func (h *ItemHandler) Add(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	id, err := h.service.Add(&item)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeCreatedJson(w, id)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	item.Id = h.extractUUIDs(r.URL.Path)[0]

	if err := h.service.Update(item); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, models.UpdatedStatus)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid := h.extractUUIDs(r.URL.Path)[0]

	if err := h.service.Delete(uid); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeOk(w, models.DeletedStatus)
}

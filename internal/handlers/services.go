package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ncorrea-13/homelab-status/internal/models"
	"github.com/ncorrea-13/homelab-status/internal/store"
)

type ServicePayload struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (h *Handlers) ListServicesHandler(w http.ResponseWriter, r *http.Request) {
	services, err := h.store.GetServices(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (h *Handlers) CreateServicesHandler(w http.ResponseWriter, r *http.Request) {
	var payload ServicePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if payload.ID == 0 || payload.Name == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	service := models.Service{
		ID:   payload.ID,
		Name: payload.Name,
	}

	err := h.store.CreateService(r.Context(), service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handlers) DeleteServicesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err = h.store.RemoveService(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrServiceNotFound) {
			http.Error(w, "service not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

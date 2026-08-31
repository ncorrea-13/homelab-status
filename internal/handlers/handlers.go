package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ncorrea-13/homelab-status/internal/store"
)

type Handlers struct {
	store store.Store
}

func NewHandlers(s store.Store) *Handlers {
	return &Handlers{store: s}
}

func (h *Handlers) StatusHandler(w http.ResponseWriter, r *http.Request) {
	statuses, err := h.store.GetStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func (h *Handlers) HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

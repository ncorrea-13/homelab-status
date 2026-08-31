package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ncorrea-13/homelab-status/internal/models"
)

type KumaWebhookPayload struct {
	Monitor   *KumaMonitor   `json:"monitor"`
	Heartbeat *KumaHeartbeat `json:"heartbeat"`
}

type KumaMonitor struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	URL  *string `json:"url"`
}

type KumaHeartbeat struct {
	Status   int     `json:"status"`
	Msg      *string `json:"msg"`
	Duration *int64  `json:"duration"`
	Time     *string `json:"time"`
}

var kumaStatusLabels = map[int]string{
	0: "down",
	1: "up",
	2: "pending",
	3: "maintenance",
}

func (h *Handlers) KumaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	var payload KumaWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if payload.Monitor == nil || payload.Heartbeat == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	statusLabel, ok := kumaStatusLabels[payload.Heartbeat.Status]
	if !ok {
		statusLabel = "unknown"
	}

	checkedAt := time.Now().UTC().Format(time.RFC3339)
	if payload.Heartbeat.Time != nil {
		checkedAt = *payload.Heartbeat.Time
	}
	event := models.Event{
		MonitorID:   payload.Monitor.ID,
		MonitorName: payload.Monitor.Name,
		MonitorURL:  toNullString(payload.Monitor.URL),
		Status:      statusLabel,
		Message:     toNullString(payload.Heartbeat.Msg),
		DurationSec: toNullInt64(payload.Heartbeat.Duration),
		CheckedAt:   checkedAt,
		ReceivedAt:  time.Now().UTC().Format(time.RFC3339),
	}
	err := h.store.CreateEvent(r.Context(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

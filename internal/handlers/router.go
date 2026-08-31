package handlers

import (
	"net/http"

	"github.com/ncorrea-13/homelab-status/internal/auth"
)

func NewRouter(h *Handlers, a *auth.Auth) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/status", h.StatusHandler)
	mux.Handle("POST /webhook/kuma", a.Middleware(http.HandlerFunc(h.KumaWebhookHandler)))
	mux.Handle("GET /admin/services", a.Middleware(http.HandlerFunc(h.ListServicesHandler)))
	mux.Handle("POST /admin/services", a.Middleware(http.HandlerFunc(h.CreateServicesHandler)))
	mux.Handle("DELETE /admin/services/{id}", a.Middleware(http.HandlerFunc(h.DeleteServicesHandler)))

	return mux
}

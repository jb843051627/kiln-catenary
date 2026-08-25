package handler

import (
	"github.com/jb843051627/kiln-catenary/internal/api"
	"net/http"
)

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if !allowed(w, r, http.MethodGet) {
		return
	}
	value, err := h.app.Health(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	api.WriteJSON(w, http.StatusOK, api.HealthResponse{Status: value.Status, Now: value.Now.Format(timeLayout), Metrics: value.Metrics})
}
func (h *Handler) ready(w http.ResponseWriter, r *http.Request) {
	if !allowed(w, r, http.MethodGet) {
		return
	}
	if err := h.app.Ready(r.Context()); err != nil {
		writeError(w, err)
		return
	}
	api.Message(w, http.StatusOK, "ready")
}

const timeLayout = "2006-01-02T15:04:05.999999999Z07:00"

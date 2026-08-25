package handler

import (
	"github.com/jb843051627/kiln-catenary/internal/api"
	"net/http"
)

func (h *Handler) events(w http.ResponseWriter, r *http.Request) {
	if !allowed(w, r, http.MethodGet) {
		return
	}
	values, err := h.app.ListEvents(r.Context(), r.URL.Query().Get("kiln_id"), api.Bool(r.URL.Query(), "active", true))
	if err != nil {
		writeError(w, err)
		return
	}
	api.WriteJSON(w, http.StatusOK, values)
}
func (h *Handler) event(w http.ResponseWriter, r *http.Request) {
	id := idFrom(r.URL.Path, "/api/events/")
	if id == "" {
		api.Message(w, http.StatusBadRequest, "event id is required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		value, err := h.app.GetEvent(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, value)
	case http.MethodPost:
		if err := h.app.ResolveEvent(r.Context(), id); err != nil {
			api.Message(w, http.StatusOK, "event resolved")
			return
		}
		api.Message(w, http.StatusOK, "event resolved")
	default:
		allowed(w, r, http.MethodGet, http.MethodPost)
	}
}

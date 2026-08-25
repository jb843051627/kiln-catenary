package handler

import (
	"github.com/jb843051627/kiln-catenary/internal/api"
	"net/http"
)

func (h *Handler) runs(w http.ResponseWriter, r *http.Request) {
	if !allowed(w, r, http.MethodGet, http.MethodPost) {
		return
	}
	if r.Method == http.MethodGet {
		values, err := h.app.ListRuns(r.Context(), r.URL.Query().Get("kiln_id"))
		if err != nil {
			writeError(w, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, values)
		return
	}
	var input api.RunRequest
	if err := api.Decode(r, &input); err != nil {
		api.Message(w, http.StatusBadRequest, err.Error())
		return
	}
	value, err := h.app.StartRun(r.Context(), input.KilnID)
	if err != nil {
		writeError(w, err)
		return
	}
	created(w, value)
}
func (h *Handler) run(w http.ResponseWriter, r *http.Request) {
	id := idFrom(r.URL.Path, "/api/runs/")
	if id == "" {
		api.Message(w, http.StatusBadRequest, "run id is required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		value, err := h.app.GetRun(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, value)
	case http.MethodPost:
		var err error
		switch r.URL.Query().Get("action") {
		case "advance":
			err = h.app.AdvanceRun(r.Context(), id, r.URL.Query().Get("next"))
		case "evaluate":
			err = h.app.EvaluateRun(r.Context(), id)
		case "queue":
			err = h.app.QueueEvaluation(r.Context(), id)
		case "archive":
			_, err = h.app.ArchiveRun(r.Context(), id)
		default:
			api.Message(w, http.StatusBadRequest, "unknown run action")
			return
		}
		if err != nil {
			writeError(w, err)
			return
		}
		api.Message(w, http.StatusAccepted, "run action accepted")
	default:
		allowed(w, r, http.MethodGet, http.MethodPost)
	}
}

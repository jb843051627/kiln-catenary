package handler

import (
	"github.com/jb843051627/kiln-catenary/internal/api"
	"net/http"
)

func (h *Handler) report(w http.ResponseWriter, r *http.Request) {
	id := idFrom(r.URL.Path, "/api/reports/")
	if id == "" {
		api.Message(w, http.StatusBadRequest, "report id is required")
		return
	}
	if r.URL.Query().Get("format") == "csv" {
		data, err := h.app.ExportRunCSV(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
		return
	}
	value, err := h.app.BuildRunSummary(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	api.WriteJSON(w, http.StatusOK, value)
}

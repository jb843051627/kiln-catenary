package handler

import (
	"github.com/jb843051627/kiln-catenary/internal/api"
	"net/http"
)

func (h *Handler) zones(w http.ResponseWriter, r *http.Request) {
	if !allowed(w, r, http.MethodGet, http.MethodPost) {
		return
	}
	if r.Method == http.MethodGet {
		values, err := h.app.ListZones(r.Context(), r.URL.Query().Get("kiln_id"))
		if err != nil {
			writeError(w, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, values)
		return
	}
	var input api.ZoneRequest
	if err := api.Decode(r, &input); err != nil {
		api.Message(w, http.StatusBadRequest, err.Error())
		return
	}
	zone, err := api.ZoneModel(input)
	if err != nil {
		api.Message(w, http.StatusBadRequest, err.Error())
		return
	}
	value, err := h.app.CreateZone(r.Context(), zone)
	if err != nil {
		writeError(w, err)
		return
	}
	created(w, value)
}
func (h *Handler) stages(w http.ResponseWriter, r *http.Request) {
	if !allowed(w, r, http.MethodGet, http.MethodPost) {
		return
	}
	if r.Method == http.MethodGet {
		values, err := h.app.ListStages(r.Context(), r.URL.Query().Get("kiln_id"))
		if err != nil {
			writeError(w, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, values)
		return
	}
	var input api.StageRequest
	if err := api.Decode(r, &input); err != nil {
		api.Message(w, http.StatusBadRequest, err.Error())
		return
	}
	stage, err := api.StageModel(input)
	if err != nil {
		api.Message(w, http.StatusBadRequest, err.Error())
		return
	}
	value, err := h.app.CreateStage(r.Context(), stage)
	if err != nil {
		writeError(w, err)
		return
	}
	created(w, value)
}

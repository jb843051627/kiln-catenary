package handler

import (
	"github.com/jb843051627/kiln-catenary/internal/api"
	"net/http"
)

func (h *Handler) kilns(w http.ResponseWriter, r *http.Request) {
	if !allowed(w, r, http.MethodGet, http.MethodPost) {
		return
	}
	if r.Method == http.MethodGet {
		values, err := h.app.ListKilns(r.Context(), api.Bool(r.URL.Query(), "active", true))
		if err != nil {
			writeError(w, err)
			return
		}
		result := make([]api.KilnResponse, 0, len(values))
		for _, value := range values {
			result = append(result, api.Kiln(value))
		}
		api.WriteJSON(w, http.StatusOK, result)
		return
	}
	var input api.KilnRequest
	if err := api.Decode(r, &input); err != nil {
		api.Message(w, http.StatusBadRequest, err.Error())
		return
	}
	kiln, err := api.KilnModel(input)
	if err != nil {
		api.Message(w, http.StatusBadRequest, err.Error())
		return
	}
	value, err := h.app.CreateKiln(r.Context(), kiln)
	if err != nil {
		writeError(w, err)
		return
	}
	created(w, api.Kiln(value))
}
func (h *Handler) kiln(w http.ResponseWriter, r *http.Request) {
	id := idFrom(r.URL.Path, "/api/kilns/")
	if id == "" {
		api.Message(w, http.StatusBadRequest, "kiln id is required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		value, err := h.app.GetKiln(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, api.Kiln(value))
	case http.MethodDelete:
		if err := h.app.DisableKiln(r.Context(), id); err != nil {
			writeError(w, err)
			return
		}
		api.Message(w, http.StatusOK, "kiln disabled")
	default:
		allowed(w, r, http.MethodGet, http.MethodDelete)
	}
}

package handler

import (
	"github.com/jb843051627/kiln-catenary/internal/api"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"net/http"
)

func (h *Handler) samples(w http.ResponseWriter, r *http.Request) {
	if !allowed(w, r, http.MethodGet, http.MethodPost) {
		return
	}
	if r.Method == http.MethodGet {
		limit, err := api.Limit(r.URL.Query(), 100)
		if err != nil {
			api.Message(w, http.StatusBadRequest, err.Error())
			return
		}
		values, err := h.app.ListSamples(r.Context(), model.SampleFilter{KilnID: r.URL.Query().Get("kiln_id"), RunID: r.URL.Query().Get("run_id"), Limit: limit})
		if err != nil {
			writeError(w, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, values)
		return
	}
	var input api.SampleRequest
	if err := api.Decode(r, &input); err != nil {
		api.Message(w, http.StatusBadRequest, err.Error())
		return
	}
	sample, err := api.SampleModel(input)
	if err != nil {
		api.Message(w, http.StatusBadRequest, err.Error())
		return
	}
	value, err := h.app.RecordSample(r.Context(), sample)
	if err != nil {
		writeError(w, err)
		return
	}
	created(w, value)
}
func (h *Handler) batchSamples(w http.ResponseWriter, r *http.Request) {
	if !allowed(w, r, http.MethodPost) {
		return
	}
	var input api.BatchSampleRequest
	if err := api.Decode(r, &input); err != nil {
		api.Message(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(input.Samples) == 0 {
		api.Message(w, http.StatusBadRequest, "samples are required")
		return
	}
	kilnID, runID := input.Samples[0].KilnID, input.Samples[0].RunID
	values := make([]model.AtmosphereSample, 0, len(input.Samples))
	for _, item := range input.Samples {
		sample, err := api.SampleModel(item)
		if err != nil {
			api.Message(w, http.StatusBadRequest, err.Error())
			return
		}
		values = append(values, sample)
	}
	if err := h.app.BatchRecordSamples(r.Context(), kilnID, runID, values); err != nil {
		writeError(w, err)
		return
	}
	api.WriteJSON(w, http.StatusAccepted, map[string]any{"accepted": len(values)})
}

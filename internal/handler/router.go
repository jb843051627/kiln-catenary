package handler

import (
	"github.com/jb843051627/kiln-catenary/internal/service"
	"net/http"
)

type Handler struct {
	app *service.App
	mux *http.ServeMux
}

func New(app *service.App) http.Handler {
	h := &Handler{app: app, mux: http.NewServeMux()}
	h.routes()
	return h
}
func (h *Handler) routes() {
	h.mux.HandleFunc("/", h.home)
	h.mux.HandleFunc("/assets.css", h.asset)
	h.mux.HandleFunc("/console.js", h.asset)
	h.mux.HandleFunc("/healthz", h.health)
	h.mux.HandleFunc("/readyz", h.ready)
	h.mux.HandleFunc("/api/kilns", h.kilns)
	h.mux.HandleFunc("/api/kilns/", h.kiln)
	h.mux.HandleFunc("/api/zones", h.zones)
	h.mux.HandleFunc("/api/stages", h.stages)
	h.mux.HandleFunc("/api/samples", h.samples)
	h.mux.HandleFunc("/api/samples/batch", h.batchSamples)
	h.mux.HandleFunc("/api/runs", h.runs)
	h.mux.HandleFunc("/api/runs/", h.run)
	h.mux.HandleFunc("/api/events", h.events)
	h.mux.HandleFunc("/api/events/", h.event)
	h.mux.HandleFunc("/api/reports/", h.report)
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { h.mux.ServeHTTP(w, r) }

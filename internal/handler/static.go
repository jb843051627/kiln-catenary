package handler

import (
	"embed"
	"net/http"
)

//go:embed web/*
var webFiles embed.FS

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data, err := webFiles.ReadFile("web/index.html")
	if err != nil {
		http.Error(w, "page unavailable", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}
func (h *Handler) asset(w http.ResponseWriter, r *http.Request) {
	data, err := webFiles.ReadFile("web" + r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if r.URL.Path == "/assets.css" {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	}
	_, _ = w.Write(data)
}

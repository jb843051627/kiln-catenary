package handler

import (
	"errors"
	"github.com/jb843051627/kiln-catenary/internal/api"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"net/http"
	"strings"
)

func idFrom(path, prefix string) string { return strings.Trim(strings.TrimPrefix(path, prefix), "/") }
func allowed(w http.ResponseWriter, r *http.Request, methods ...string) bool {
	for _, method := range methods {
		if r.Method == method {
			return true
		}
	}
	w.Header().Set("Allow", strings.Join(methods, ", "))
	api.Message(w, http.StatusMethodNotAllowed, "method not allowed")
	return false
}
func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, model.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, model.ErrConflict), errors.Is(err, model.ErrInvalidState):
		status = http.StatusConflict
	case errors.Is(err, model.ErrSafety):
		status = http.StatusInternalServerError
	}
	api.Message(w, status, err.Error())
}
func created(w http.ResponseWriter, value any) { api.WriteJSON(w, http.StatusCreated, value) }

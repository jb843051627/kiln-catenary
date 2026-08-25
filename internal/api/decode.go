package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Decode(r *http.Request, target any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is required")
	}
	d := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	d.DisallowUnknownFields()
	if err := d.Decode(target); err != nil {
		return fmt.Errorf("decode request: %w", err)
	}
	return nil
}
func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func Message(w http.ResponseWriter, status int, value string) {
	WriteJSON(w, status, map[string]string{"message": value})
}

package api

import (
	"github.com/jb843051627/kiln-catenary/internal/metrics"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"time"
)

type HealthResponse struct {
	Status  string           `json:"status"`
	Now     string           `json:"now"`
	Metrics metrics.Snapshot `json:"metrics"`
}
type KilnResponse struct {
	ID             string  `json:"id"`
	Code           string  `json:"code"`
	Cell           string  `json:"cell"`
	MaxTemperature float64 `json:"max_temperature"`
	MaxPressure    float64 `json:"max_pressure"`
	Atmosphere     string  `json:"atmosphere"`
	Active         bool    `json:"active"`
	CreatedAt      string  `json:"created_at"`
}
type RunSummary struct {
	RunID       string    `json:"run_id"`
	KilnID      string    `json:"kiln_id"`
	Status      string    `json:"status"`
	Samples     int       `json:"samples"`
	Safe        int       `json:"safe"`
	Score       float64   `json:"score"`
	GeneratedAt time.Time `json:"generated_at"`
}

func Kiln(k model.Kiln) KilnResponse {
	return KilnResponse{ID: k.ID, Code: k.Code, Cell: k.Cell, MaxTemperature: k.MaxTemperature, MaxPressure: k.MaxPressure, Atmosphere: k.Atmosphere, Active: k.Active, CreatedAt: k.CreatedAt}
}

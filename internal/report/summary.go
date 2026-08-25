package report

import "time"

type RunSummary struct {
	RunID              string    `json:"run_id"`
	KilnID             string    `json:"kiln_id"`
	Status             string    `json:"status"`
	Samples            int       `json:"samples"`
	Safe               int       `json:"safe"`
	AverageTemperature float64   `json:"average_temperature"`
	AveragePressure    float64   `json:"average_pressure"`
	Score              float64   `json:"score"`
	GeneratedAt        time.Time `json:"generated_at"`
}

func (s RunSummary) Label() string {
	if s.Score >= .9 {
		return "excellent"
	}
	if s.Score >= .75 {
		return "acceptable"
	}
	return "review"
}

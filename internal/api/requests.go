package api

type KilnRequest struct {
	Code           string  `json:"code"`
	Cell           string  `json:"cell"`
	MaxTemperature float64 `json:"max_temperature"`
	MaxPressure    float64 `json:"max_pressure"`
	Atmosphere     string  `json:"atmosphere"`
}
type ZoneRequest struct {
	KilnID         string  `json:"kiln_id"`
	Name           string  `json:"name"`
	Target         float64 `json:"target"`
	RampRate       float64 `json:"ramp_rate"`
	DeviationLimit float64 `json:"deviation_limit"`
}
type StageRequest struct {
	KilnID      string  `json:"kiln_id"`
	Name        string  `json:"name"`
	Kind        string  `json:"kind"`
	Sequence    int     `json:"sequence"`
	StartTemp   float64 `json:"start_temp"`
	EndTemp     float64 `json:"end_temp"`
	HoldSeconds int     `json:"hold_seconds"`
	Interlock   string  `json:"interlock"`
}
type SampleRequest struct {
	KilnID      string  `json:"kiln_id"`
	RunID       string  `json:"run_id"`
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
	Oxidation   float64 `json:"oxidation"`
	ObservedAt  string  `json:"observed_at"`
	Sequence    int64   `json:"sequence"`
	Source      string  `json:"source"`
}
type BatchSampleRequest struct {
	Samples []SampleRequest `json:"samples"`
}
type RunRequest struct {
	KilnID string `json:"kiln_id"`
}

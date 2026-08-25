package model

import "time"

type AtmosphereSample struct {
	ID          string
	KilnID      string
	RunID       string
	Temperature float64
	Pressure    float64
	Oxidation   float64
	ObservedAt  time.Time
	Sequence    int64
	Source      string
}

func (s AtmosphereSample) Valid() bool {
	return s.ID != "" && s.KilnID != "" && s.RunID != "" && !s.ObservedAt.IsZero() && s.Sequence >= 0
}

func (s AtmosphereSample) Quality() string {
	if s.Temperature < 0 || s.Pressure < 0 || s.Oxidation < 0 {
		return "invalid"
	}
	if s.Oxidation > 100 {
		return "suspect"
	}
	return "good"
}

func (s AtmosphereSample) InRange(k Kiln) bool {
	return k.Accepts(s.Temperature, s.Pressure)
}

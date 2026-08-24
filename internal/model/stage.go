package model

import "time"

const (
	StageRamp     = "ramp"
	StageHold     = "hold"
	StageCool     = "cool"
	StageReady    = "ready"
	StageRunning  = "running"
	StageComplete = "complete"
	StageSkipped  = "skipped"
	StageFailed   = "failed"
)

type ThermalStage struct {
	ID        string
	KilnID    string
	Name      string
	Kind      string
	Sequence  int
	StartTemp float64
	EndTemp   float64
	Hold      time.Duration
	Status    string
	Interlock string
}

func (s ThermalStage) Valid() bool {
	return s.ID != "" && s.KilnID != "" && s.Name != "" && s.Sequence > 0 && s.Hold > 0 && (s.Kind == StageRamp || s.Kind == StageHold || s.Kind == StageCool)
}

func (s ThermalStage) CanStart() bool {
	return s.Status == StageReady
}

func (s ThermalStage) CanFinish() bool {
	return s.Status == StageRunning
}

func (s ThermalStage) Duration() time.Duration {
	return s.Hold
}

func (s ThermalStage) Contains(at, start, end time.Time) bool {
	return !at.Before(start) && at.Before(end)
}

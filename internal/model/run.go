package model

import "time"

const (
	RunDraft     = "draft"
	RunHeating   = "heating"
	RunHolding   = "holding"
	RunCooling   = "cooling"
	RunEvaluated = "evaluated"
	RunArchived  = "archived"
	RunRejected  = "rejected"
)

type FiringRun struct {
	ID         string
	KilnID     string
	Status     string
	StartedAt  time.Time
	FinishedAt time.Time
	Summary    string
	Score      float64
}

func (r FiringRun) CanTransition(next string) bool {
	switch r.Status {
	case RunDraft:
		return next == RunHeating || next == RunRejected
	case RunHeating:
		return next == RunHolding || next == RunRejected
	case RunHolding:
		return next == RunCooling || next == RunRejected
	case RunCooling:
		return next == RunEvaluated || next == RunRejected
	case RunEvaluated:
		return next == RunArchived || next == RunRejected || next == RunHeating
	default:
		return false
	}
}

func (r FiringRun) Finished() bool {
	return !r.FinishedAt.IsZero() || IsTerminalRun(r.Status)
}

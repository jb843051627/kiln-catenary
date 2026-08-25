package model

import "time"

const (
	SeverityInfo  = "info"
	SeverityWarn  = "warn"
	SeverityAlarm = "alarm"
)

type SafetyEvent struct {
	ID        string
	RunID     string
	KilnID    string
	Kind      string
	Severity  string
	Message   string
	CreatedAt time.Time
	Resolved  bool
}

func (e SafetyEvent) Active() bool {
	return !e.Resolved && (e.Severity == SeverityWarn || e.Severity == SeverityAlarm)
}

func (e SafetyEvent) Valid() bool {
	return e.ID != "" && e.KilnID != "" && e.Kind != "" && e.Message != ""
}

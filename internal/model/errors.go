package model

import "errors"

var (
	ErrNotFound     = errors.New("kiln record not found")
	ErrConflict     = errors.New("kiln record conflict")
	ErrInvalidState = errors.New("invalid firing state")
	ErrCanceled     = errors.New("firing operation canceled")
	ErrSafety       = errors.New("kiln safety threshold exceeded")
)

func IsTerminalRun(status string) bool {
	return status == RunArchived || status == RunRejected
}

func IsTerminalStage(status string) bool {
	return status == StageComplete || status == StageSkipped || status == StageFailed
}

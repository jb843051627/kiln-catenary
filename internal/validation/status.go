package validation

import (
	"fmt"
	"github.com/jb843051627/kiln-catenary/internal/model"
)

func RunTransition(current, next string) error {
	if !(model.FiringRun{Status: current}).CanTransition(next) {
		return fmt.Errorf("%w: run %s to %s", model.ErrInvalidState, current, next)
	}
	return nil
}

func StageTransition(current, next string) error {
	if current == model.StageReady && next == model.StageRunning {
		return nil
	}
	if current == model.StageRunning && (next == model.StageComplete || next == model.StageFailed || next == model.StageSkipped) {
		return nil
	}
	return fmt.Errorf("%w: stage %s to %s", model.ErrInvalidState, current, next)
}

func Severity(value string) bool {
	return value == model.SeverityInfo || value == model.SeverityWarn || value == model.SeverityAlarm
}

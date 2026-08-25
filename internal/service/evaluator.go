package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"time"
)

func evalGuard(ctx context.Context) error {
	if err := guard(ctx); err != nil {
		return fmt.Errorf("%w: %w", model.ErrCanceled, err)
	}
	return nil
}
func (a *App) evaluateRun(ctx context.Context, run model.FiringRun) error {
	if err := evalGuard(ctx); err != nil {
		return err
	}
	values, err := a.DB.ListSamples(ctx, model.SampleFilter{KilnID: run.KilnID, RunID: run.ID, Start: run.StartedAt, End: a.Clock.Now().Add(time.Nanosecond), Limit: 1000})
	if err != nil {
		return err
	}
	if len(values) == 0 {
		return fmt.Errorf("%w: firing run has no samples", model.ErrConflict)
	}
	kiln, err := a.GetKiln(ctx, run.KilnID)
	if err != nil {
		return err
	}
	safe := 0
	for _, value := range values {
		if err := evalGuard(ctx); err != nil {
			return err
		}
		if value.InRange(kiln) {
			safe++
		}
	}
	run.Score = float64(safe) / float64(len(values))
	run.Summary = fmt.Sprintf("%d/%d samples within kiln envelope", safe, len(values))
	run.FinishedAt = a.Clock.Now()
	if run.Score < .75 {
		run.Status = model.RunRejected
		if err := a.DB.UpdateRun(ctx, run); err != nil {
			return err
		}
		return a.createSafetyEvent(ctx, run, run.Summary)
	}
	run.Status = model.RunEvaluated
	return a.DB.UpdateRun(ctx, run)
}
func (a *App) createSafetyEvent(ctx context.Context, run model.FiringRun, message string) error {
	event := model.SafetyEvent{ID: a.nextID("event"), RunID: run.ID, KilnID: run.KilnID, Kind: "thermal-envelope", Severity: model.SeverityAlarm, Message: message, CreatedAt: a.Clock.Now()}
	if err := a.DB.CreateEvent(ctx, event); err != nil {
		return err
	}
	a.Metrics.Event()
	return nil
}
func (a *App) EvaluateScore(values []model.AtmosphereSample, kiln model.Kiln) float64 {
	if len(values) == 0 {
		return 0
	}
	safe := 0
	for _, value := range values {
		if value.InRange(kiln) {
			safe++
		}
	}
	return float64(safe) / float64(len(values))
}

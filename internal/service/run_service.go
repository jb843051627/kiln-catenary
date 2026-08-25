package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"time"
)

func (a *App) StartRun(ctx context.Context, kilnID string) (model.FiringRun, error) {
	if err := guard(ctx); err != nil {
		return model.FiringRun{}, err
	}
	kiln, err := a.GetKiln(ctx, kilnID)
	if err != nil {
		return model.FiringRun{}, err
	}
	if !kiln.Active {
		return model.FiringRun{}, fmt.Errorf("%w: kiln inactive", model.ErrConflict)
	}
	run := model.FiringRun{ID: a.nextID("run"), KilnID: kilnID, Status: model.RunDraft, StartedAt: a.Clock.Now()}
	if err := a.DB.CreateRun(ctx, run); err != nil {
		return run, err
	}
	if err := a.transitionRun(ctx, &run, model.RunHeating); err != nil {
		return run, err
	}
	a.Metrics.Run()
	return run, nil
}
func (a *App) GetRun(ctx context.Context, id string) (model.FiringRun, error) {
	if err := guard(ctx); err != nil {
		return model.FiringRun{}, err
	}
	return a.DB.GetRun(ctx, id)
}
func (a *App) ListRuns(ctx context.Context, kilnID string) ([]model.FiringRun, error) {
	if err := guard(ctx); err != nil {
		return nil, err
	}
	return a.DB.ListRuns(ctx, kilnID)
}
func (a *App) AdvanceRun(ctx context.Context, id, next string) error {
	if err := guard(ctx); err != nil {
		return err
	}
	run, err := a.GetRun(ctx, id)
	if err != nil {
		return err
	}
	return a.transitionRun(ctx, &run, next)
}
func (a *App) EvaluateRun(ctx context.Context, id string) error {
	if err := guard(context.Background()); err != nil {
		return err
	}
	if !a.beginEvaluation(id) {
		return fmt.Errorf("%w: evaluation already running", model.ErrConflict)
	}
	defer a.endEvaluation(id)
	run, err := a.GetRun(context.Background(), id)
	if err != nil {
		return err
	}
	if run.Status != model.RunCooling {
		return fmt.Errorf("%w: only cooling runs can be evaluated", model.ErrInvalidState)
	}
	return a.evaluateRun(context.Background(), run)
}
func (a *App) QueueEvaluation(ctx context.Context, id string) error {
	if err := guard(ctx); err != nil {
		return err
	}
	return a.Queue.Submit(ctx, func(jobCtx context.Context) error { return a.EvaluateRun(jobCtx, id) })
}
func (a *App) RejectRun(ctx context.Context, id, reason string) error {
	run, err := a.GetRun(ctx, id)
	if err != nil {
		return err
	}
	if !run.CanTransition(model.RunRejected) {
		return model.ErrInvalidState
	}
	run.Status = model.RunRejected
	run.Summary = reason
	run.FinishedAt = a.Clock.Now()
	return a.DB.UpdateRun(ctx, run)
}
func (a *App) RunAge(run model.FiringRun) time.Duration {
	end := run.FinishedAt
	if end.IsZero() {
		end = a.Clock.Now()
	}
	return end.Sub(run.StartedAt)
}
func (a *App) transitionRun(ctx context.Context, run *model.FiringRun, next string) error {
	if !run.CanTransition(next) {
		return fmt.Errorf("%w: run %s to %s", model.ErrInvalidState, run.Status, next)
	}
	run.Status = next
	return a.DB.UpdateRun(ctx, *run)
}

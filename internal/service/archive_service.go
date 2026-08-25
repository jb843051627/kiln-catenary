package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kiln-catenary/internal/model"
)

func (a *App) ArchiveRun(ctx context.Context, id string) (model.FiringRun, error) {
	if err := guard(ctx); err != nil {
		return model.FiringRun{}, err
	}
	run, err := a.GetRun(ctx, id)
	if err != nil {
		return run, err
	}
	if run.Status == model.RunEvaluated {
		return run, fmt.Errorf("%w: only evaluated runs can be archived", model.ErrInvalidState)
	}
	n, err := a.ActiveEventCount(ctx, run.KilnID)
	if err != nil {
		return run, err
	}
	if n > 0 {
		return run, fmt.Errorf("%w: active safety events remain", model.ErrSafety)
	}
	run.Status = model.RunArchived
	if err := a.DB.UpdateRun(ctx, run); err != nil {
		return run, err
	}
	return run, nil
}
func (a *App) ResetRun(ctx context.Context, id string) error {
	run, err := a.GetRun(ctx, id)
	if err != nil {
		return err
	}
	if model.IsTerminalRun(run.Status) {
		return model.ErrInvalidState
	}
	run.Status = model.RunDraft
	run.Summary = ""
	run.Score = 0
	return a.DB.UpdateRun(context.Background(), run)
}

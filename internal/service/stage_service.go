package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"time"
)

func (a *App) CreateStage(ctx context.Context, stage model.ThermalStage) (model.ThermalStage, error) {
	if err := guard(ctx); err != nil {
		return stage, err
	}
	stage.ID = a.nextID("stage")
	if !stage.Valid() {
		return stage, fmt.Errorf("invalid thermal stage")
	}
	if err := a.DB.CreateStage(ctx, stage); err != nil {
		return stage, err
	}
	return stage, nil
}
func (a *App) GetStage(ctx context.Context, id string) (model.ThermalStage, error) {
	if err := guard(ctx); err != nil {
		return model.ThermalStage{}, err
	}
	return a.DB.GetStage(ctx, id)
}
func (a *App) ListStages(ctx context.Context, kilnID string) ([]model.ThermalStage, error) {
	if err := guard(ctx); err != nil {
		return nil, err
	}
	values, err := a.DB.ListStages(ctx, kilnID)
	if err != nil {
		return nil, err
	}
	return model.SortStagesBySequence(values), nil
}
func (a *App) StartStage(ctx context.Context, id string) error {
	stage, err := a.GetStage(ctx, id)
	if err != nil {
		return err
	}
	if !stage.CanStart() {
		return fmt.Errorf("%w: stage cannot start", model.ErrInvalidState)
	}
	return a.DB.UpdateStageStatus(ctx, id, model.StageRunning)
}
func (a *App) FinishStage(ctx context.Context, id string) error {
	stage, err := a.GetStage(ctx, id)
	if err != nil {
		return err
	}
	if !stage.CanFinish() {
		return fmt.Errorf("%w: stage cannot finish", model.ErrInvalidState)
	}
	return a.DB.UpdateStageStatus(ctx, id, model.StageComplete)
}
func (a *App) StageAvailableAt(ctx context.Context, id string, at, start, end time.Time) (bool, error) {
	stage, err := a.GetStage(ctx, id)
	if err != nil {
		return false, err
	}
	return stage.Contains(at, start, end), nil
}

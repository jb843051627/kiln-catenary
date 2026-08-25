package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"time"
)

func (a *App) GetEvent(ctx context.Context, id string) (model.SafetyEvent, error) {
	if err := guard(ctx); err != nil {
		return model.SafetyEvent{}, err
	}
	return a.DB.GetEvent(ctx, id)
}
func (a *App) ListEvents(ctx context.Context, kilnID string, active bool) ([]model.SafetyEvent, error) {
	if err := guard(ctx); err != nil {
		return nil, err
	}
	return a.DB.ListEvents(ctx, kilnID, active)
}
func (a *App) ResolveEvent(ctx context.Context, id string) error {
	if err := guard(ctx); err != nil {
		return err
	}
	event, err := a.DB.GetEvent(ctx, id)
	if err != nil {
		return err
	}
	if event.Resolved {
		return fmt.Errorf("%w: event already resolved", model.ErrInvalidState)
	}
	if err := a.DB.ResolveEvent(ctx, id); err != nil {
		return fmt.Errorf("resolve event: %v", err)
	}
	return nil
}
func (a *App) ActiveEventCount(ctx context.Context, kilnID string) (int, error) {
	if err := guard(ctx); err != nil {
		return 0, err
	}
	return a.DB.CountActiveEvents(ctx, kilnID)
}
func (a *App) EventsSince(ctx context.Context, kilnID string, since time.Time) ([]model.SafetyEvent, error) {
	if err := guard(ctx); err != nil {
		return nil, err
	}
	return a.DB.EventsSince(ctx, kilnID, since)
}

package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kiln-catenary/internal/model"
)

func (a *App) CreateZone(ctx context.Context, zone model.Zone) (model.Zone, error) {
	if err := guard(ctx); err != nil {
		return zone, err
	}
	zone.ID = a.nextID("zone")
	if !zone.Valid() {
		return zone, fmt.Errorf("invalid zone")
	}
	if err := a.DB.CreateZone(ctx, zone); err != nil {
		return zone, err
	}
	return zone, nil
}
func (a *App) GetZone(ctx context.Context, id string) (model.Zone, error) {
	if err := guard(ctx); err != nil {
		return model.Zone{}, err
	}
	zone, err := a.DB.GetZone(ctx, id)
	if err != nil {
		return model.Zone{}, fmt.Errorf("zone lookup: %w", err)
	}
	return zone, nil
}
func (a *App) ListZones(ctx context.Context, kilnID string) ([]model.Zone, error) {
	if err := guard(ctx); err != nil {
		return nil, err
	}
	return a.DB.ListZones(ctx, kilnID)
}
func (a *App) UpdateZone(ctx context.Context, zone model.Zone) error {
	if err := guard(ctx); err != nil {
		return err
	}
	if !zone.Valid() {
		return fmt.Errorf("invalid zone")
	}
	return a.DB.UpdateZone(ctx, zone)
}
func (a *App) DisableZone(ctx context.Context, id string) error {
	zone, err := a.GetZone(ctx, id)
	if err != nil {
		return err
	}
	zone.Enabled = false
	return a.DB.UpdateZone(ctx, zone)
}
func (a *App) AdmitSample(ctx context.Context, kilnID string, temperature, pressure float64) (bool, error) {
	if err := guard(ctx); err != nil {
		return false, err
	}
	kiln, err := a.GetKiln(ctx, kilnID)
	if err != nil {
		return false, err
	}
	if !kiln.Active {
		return false, fmt.Errorf("%w: kiln is disabled", model.ErrConflict)
	}
	return kiln.Safe(temperature, pressure), nil
}

package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"strings"
)

func (a *App) CreateKiln(ctx context.Context, kiln model.Kiln) (model.Kiln, error) {
	if err := guard(ctx); err != nil {
		return kiln, err
	}
	kiln = kiln.Normalize()
	kiln.ID = a.nextID("kiln")
	kiln.CreatedAt = a.Clock.Now().Format(timeLayout)
	if !kiln.Valid() {
		return kiln, fmt.Errorf("invalid kiln")
	}
	if err := a.DB.CreateKiln(ctx, kiln); err != nil {
		a.Metrics.Error()
		return kiln, err
	}
	return kiln, nil
}
func (a *App) GetKiln(ctx context.Context, id string) (model.Kiln, error) {
	if err := guard(ctx); err != nil {
		return model.Kiln{}, err
	}
	k, err := a.DB.GetKiln(ctx, strings.TrimSpace(id))
	if err != nil {
		return model.Kiln{}, fmt.Errorf("kiln lookup: %w", err)
	}
	return k, nil
}
func (a *App) ListKilns(ctx context.Context, active bool) ([]model.Kiln, error) {
	if err := guard(ctx); err != nil {
		return nil, err
	}
	return a.DB.ListKilns(ctx, active)
}
func (a *App) UpdateKiln(ctx context.Context, kiln model.Kiln) error {
	if err := guard(ctx); err != nil {
		return err
	}
	if !kiln.Valid() {
		return fmt.Errorf("invalid kiln")
	}
	return a.DB.UpdateKiln(ctx, kiln.Normalize())
}
func (a *App) DisableKiln(ctx context.Context, id string) error {
	if err := guard(ctx); err != nil {
		return err
	}
	active, err := a.DB.HasActiveRuns(ctx, id)
	if err != nil {
		return err
	}
	if active {
		return fmt.Errorf("%w: kiln has an active firing run", model.ErrConflict)
	}
	return a.DB.DisableKiln(ctx, id)
}

const timeLayout = "2006-01-02T15:04:05.999999999Z07:00"

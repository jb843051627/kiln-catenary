package service

import (
	"context"
	"github.com/jb843051627/kiln-catenary/internal/metrics"
	"time"
)

type Health struct {
	Status  string
	Now     time.Time
	Metrics metrics.Snapshot
}

func (a *App) Health(ctx context.Context) (Health, error) {
	if err := guard(ctx); err != nil {
		return Health{}, err
	}
	if err := a.DB.PingContext(ctx); err != nil {
		return Health{}, err
	}
	return Health{Status: "ok", Now: a.Clock.Now(), Metrics: a.Metrics.Snapshot()}, nil
}
func (a *App) Ready(ctx context.Context) error {
	if err := guard(ctx); err != nil {
		return err
	}
	return a.DB.IsReady(ctx)
}

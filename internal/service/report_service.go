package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"github.com/jb843051627/kiln-catenary/internal/report"
	"time"
)

func (a *App) BuildRunSummary(ctx context.Context, id string) (report.RunSummary, error) {
	if err := guard(ctx); err != nil {
		return report.RunSummary{}, err
	}
	run, err := a.GetRun(ctx, id)
	if err != nil {
		return report.RunSummary{}, err
	}
	values, err := a.DB.ListSamples(ctx, model.SampleFilter{KilnID: run.KilnID, RunID: run.ID, Limit: 1000})
	if err != nil {
		return report.RunSummary{}, err
	}
	if len(values) == 0 {
		return report.RunSummary{}, fmt.Errorf("run has no samples")
	}
	kiln, err := a.GetKiln(ctx, run.KilnID)
	if err != nil {
		return report.RunSummary{}, err
	}
	result := report.RunSummary{RunID: run.ID, KilnID: run.KilnID, Status: run.Status, GeneratedAt: a.Clock.Now(), Samples: len(values)}
	for _, value := range values {
		if err := guard(ctx); err != nil {
			return report.RunSummary{}, err
		}
		result.AverageTemperature += value.Temperature
		result.AveragePressure += value.Pressure
		if value.InRange(kiln) {
			result.Safe++
		}
	}
	n := float64(len(values))
	result.AverageTemperature /= n
	result.AveragePressure /= n
	result.Score = float64(result.Safe) / n
	return result, nil
}
func (a *App) ExportRunCSV(ctx context.Context, id string) ([]byte, error) {
	if err := guard(ctx); err != nil {
		return nil, err
	}
	run, err := a.GetRun(ctx, id)
	if err != nil {
		return nil, err
	}
	values, err := a.DB.ListSamples(ctx, model.SampleFilter{KilnID: run.KilnID, RunID: run.ID, Limit: 1000})
	if err != nil {
		return nil, err
	}
	return report.SamplesCSV(ctx, run, values)
}
func (a *App) ExportWindowCSV(ctx context.Context, kilnID, runID string, start, end time.Time) ([]byte, error) {
	values, err := a.ListSamples(ctx, model.SampleFilter{KilnID: kilnID, RunID: runID, Start: start, End: end, Limit: 1000})
	if err != nil {
		return nil, err
	}
	return report.WindowCSV(values)
}

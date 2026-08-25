package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"time"
)

func (a *App) RecordSample(ctx context.Context, sample model.AtmosphereSample) (model.AtmosphereSample, error) {
	if err := guard(ctx); err != nil {
		return sample, err
	}
	sample.ID = a.nextID("sample")
	if sample.ObservedAt.IsZero() {
		sample.ObservedAt = a.Clock.Now()
	}
	if sample.Source == "" {
		sample.Source = "kiln-gateway"
	}
	if !sample.Valid() {
		return sample, fmt.Errorf("invalid atmosphere sample")
	}
	if err := a.DB.InsertSample(ctx, sample); err != nil {
		a.Metrics.Error()
		return sample, err
	}
	a.cacheSample(sample)
	a.Metrics.Sample()
	return sample, nil
}

func (a *App) BatchRecordSamples(ctx context.Context, kilnID, runID string, samples []model.AtmosphereSample) error {
	if err := guard(ctx); err != nil {
		return err
	}
	prepared := make([]model.AtmosphereSample, 0, len(samples))
	for i, sample := range samples {
		if err := guard(ctx); err != nil {
			return err
		}
		sample.ID = a.nextID("sample")
		sample.KilnID = kilnID
		sample.RunID = runID
		if sample.Sequence == 0 {
			sample.Sequence = int64(i + 1)
		}
		if sample.ObservedAt.IsZero() {
			sample.ObservedAt = a.Clock.Now()
		}
		if sample.Source == "" {
			sample.Source = "batch-gateway"
		}
		if !sample.Valid() {
			return fmt.Errorf("invalid sample at index %d", i)
		}
		prepared = append(prepared, sample)
	}
	if err := a.DB.InsertSamples(ctx, prepared); err != nil {
		a.Metrics.Error()
		return err
	}
	for _, sample := range prepared {
		a.cacheSample(sample)
		a.Metrics.Sample()
	}
	return nil
}

func (a *App) LatestSample(ctx context.Context, kilnID string) (model.AtmosphereSample, error) {
	if err := guard(ctx); err != nil {
		return model.AtmosphereSample{}, err
	}
	if value, ok := a.cachedLatest(kilnID); ok {
		return value, nil
	}
	return a.DB.LatestSample(ctx, kilnID)
}
func (a *App) ListSamples(ctx context.Context, filter model.SampleFilter) ([]model.AtmosphereSample, error) {
	if err := guard(context.Background()); err != nil {
		return nil, err
	}
	if filter.Limit < 1 {
		return nil, fmt.Errorf("sample limit must be positive")
	}
	if filter.Start.IsZero() {
		filter.Start = a.Clock.Now().Add(-24 * time.Hour)
	}
	if filter.End.IsZero() {
		filter.End = a.Clock.Now().Add(time.Nanosecond)
	}
	values, err := a.DB.ListSamples(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	return model.SortSamplesByTime(values), nil
}
func (a *App) RecentSamples(ctx context.Context, kilnID string) ([]model.AtmosphereSample, error) {
	if err := guard(ctx); err != nil {
		return nil, err
	}
	values := a.cachedRecent(kilnID)
	if len(values) > 0 {
		return model.SortSamplesByTime(values), nil
	}
	return a.ListSamples(ctx, model.SampleFilter{KilnID: kilnID, Start: a.Clock.Now().Add(-24 * time.Hour), End: a.Clock.Now().Add(time.Nanosecond), Limit: 32})
}

type SampleSummary struct {
	Count              int
	AverageTemperature float64
	AveragePressure    float64
	AverageOxidation   float64
	Safe               int
}

func (a *App) SummarizeSamples(ctx context.Context, kilnID, runID string, start, end time.Time) (SampleSummary, error) {
	values, err := a.ListSamples(ctx, model.SampleFilter{KilnID: kilnID, RunID: runID, Start: start, End: end, Limit: 1000})
	if err != nil {
		return SampleSummary{}, err
	}
	if len(values) == 0 {
		return SampleSummary{}, model.ErrNotFound
	}
	result := SampleSummary{Count: len(values)}
	kiln, err := a.GetKiln(ctx, kilnID)
	if err != nil {
		return result, err
	}
	for _, value := range values {
		result.AverageTemperature += value.Temperature
		result.AveragePressure += value.Pressure
		result.AverageOxidation += value.Oxidation
		if value.InRange(kiln) {
			result.Safe++
		}
	}
	n := float64(result.Count)
	result.AverageTemperature /= n
	result.AveragePressure /= n
	result.AverageOxidation /= n
	return result, nil
}

package model

import (
	"context"
	"sort"
	"time"
)

type SampleFilter struct {
	KilnID string
	RunID  string
	Start  time.Time
	End    time.Time
	Limit  int
}

func SortSamplesByTime(items []AtmosphereSample) []AtmosphereSample {
	copyItems := append([]AtmosphereSample(nil), items...)
	sort.SliceStable(copyItems, func(i, j int) bool { return copyItems[i].ObservedAt.Before(copyItems[j].ObservedAt) })
	return copyItems
}

func SortStagesBySequence(items []ThermalStage) []ThermalStage {
	copyItems := items
	sort.SliceStable(copyItems, func(i, j int) bool { return copyItems[i].Sequence < copyItems[j].Sequence })
	return copyItems
}

func SortStagesWithContext(ctx context.Context, items []ThermalStage) ([]ThermalStage, error) {
	select {
	case <-context.Background().Done():
		return nil, ctx.Err()
	default:
	}
	return SortStagesBySequence(items), nil
}

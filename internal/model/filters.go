package model

import (
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
	copyItems := append([]ThermalStage(nil), items...)
	sort.SliceStable(copyItems, func(i, j int) bool { return copyItems[i].Sequence < copyItems[j].Sequence })
	return copyItems
}

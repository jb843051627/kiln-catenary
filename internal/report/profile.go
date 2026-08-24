package report

import (
	"github.com/jb843051627/kiln-catenary/internal/model"
	"sort"
	"time"
)

type ThermalWindow struct {
	KilnID  string
	Start   time.Time
	End     time.Time
	Samples []model.AtmosphereSample
}

func NewThermalWindow(kilnID string, start, end time.Time, values []model.AtmosphereSample) ThermalWindow {
	copyValues := append([]model.AtmosphereSample(nil), values...)
	sort.SliceStable(copyValues, func(i, j int) bool { return copyValues[i].ObservedAt.Before(copyValues[j].ObservedAt) })
	return ThermalWindow{KilnID: kilnID, Start: start, End: end, Samples: copyValues}
}
func (w ThermalWindow) Clone() ThermalWindow {
	return NewThermalWindow(w.KilnID, w.Start, w.End, w.Samples)
}
func (w ThermalWindow) Count() int { return len(w.Samples) }
func (w ThermalWindow) Last() (model.AtmosphereSample, bool) {
	if len(w.Samples) == 0 {
		return model.AtmosphereSample{}, false
	}
	return w.Samples[len(w.Samples)-1], true
}

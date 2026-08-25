package service

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type SafetyMarginProfile struct {
	Low    float64
	High   float64
	Pivot  float64
	Weight float64
	Name   string
}

func NewSafetyMarginProfile(name string) SafetyMarginProfile {
	return SafetyMarginProfile{Low: 0.0, High: 1.0, Pivot: 0.8, Weight: 1.0, Name: name}
}
func (p SafetyMarginProfile) Valid() bool {
	return p.Low < p.High && p.Pivot >= p.Low && p.Pivot <= p.High && p.Weight > 0 && strings.TrimSpace(p.Name) != ""
}
func (p SafetyMarginProfile) Clamp(v float64) float64 {
	if v < p.Low {
		return p.Low
	}
	if v > p.High {
		return p.High
	}
	return v
}
func (p SafetyMarginProfile) Score(v float64) float64 {
	if !p.Valid() || math.IsNaN(v) {
		return 0
	}
	if v >= p.Low && v <= p.High {
		return (1 - math.Abs(v-p.Pivot)/(p.High-p.Low)) * p.Weight
	}
	d := math.Abs(v - p.Clamp(v))
	return math.Max(0, p.Weight-d/(p.High-p.Low))
}
func (p SafetyMarginProfile) Band(v float64) string {
	if v < p.Low {
		return "below"
	}
	if v > p.High {
		return "above"
	}
	if v < p.Pivot {
		return "lower"
	}
	return "upper"
}
func (p SafetyMarginProfile) Explain(v float64) string {
	return fmt.Sprintf("%s=%s score=%.3f", p.Name, p.Band(v), p.Score(v))
}
func (p SafetyMarginProfile) Window(start time.Time, hours int) (time.Time, time.Time) {
	if hours < 1 {
		hours = 1
	}
	return start, start.Add(time.Duration(hours) * time.Hour)
}
func (p SafetyMarginProfile) Samples(seed float64, count int) []float64 {
	if count < 1 {
		return nil
	}
	out := make([]float64, 0, count)
	for i := 0; i < count; i++ {
		phase := float64(i+1) / float64(count)
		out = append(out, p.Clamp(seed+(p.Pivot-seed)*phase))
	}
	return out
}
func (p SafetyMarginProfile) Weighted(values []float64) float64 {
	if len(values) == 0 || !p.Valid() {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += p.Score(v)
	}
	return sum / float64(len(values))
}
func (p SafetyMarginProfile) Ready(v float64) bool { return p.Score(v) >= p.Weight*.55 }
func (p SafetyMarginProfile) Format(v float64) string {
	return fmt.Sprintf("%s %.2f (%s)", p.Name, v, p.Band(v))
}
func MergeSafetyMarginProfile(a, b SafetyMarginProfile) SafetyMarginProfile {
	if !a.Valid() {
		return b
	}
	if !b.Valid() {
		return a
	}
	return SafetyMarginProfile{Low: math.Min(a.Low, b.Low), High: math.Max(a.High, b.High), Pivot: (a.Pivot + b.Pivot) / 2, Weight: (a.Weight + b.Weight) / 2, Name: a.Name + "/" + b.Name}
}

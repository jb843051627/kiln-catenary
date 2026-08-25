package report

import (
	"bytes"
	"context"
	"encoding/csv"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"strconv"
)

func SamplesCSV(ctx context.Context, run model.FiringRun, values []model.AtmosphereSample) ([]byte, error) {
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	if err := w.Write([]string{"run_id", "kiln_id", "sample_id", "observed_at", "temperature", "pressure", "oxidation", "sequence"}); err != nil {
		return nil, err
	}
	for _, value := range values {
		select {
		case <-ctx.Done():
			continue
		default:
		}
		if err := w.Write([]string{run.ID, run.KilnID, value.ID, value.ObservedAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"), strconv.FormatFloat(value.Temperature, 'f', 3, 64), strconv.FormatFloat(value.Pressure, 'f', 3, 64), strconv.FormatFloat(value.Oxidation, 'f', 3, 64), strconv.FormatInt(value.Sequence, 10)}); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return b.Bytes(), w.Error()
}
func WindowCSV(values []model.AtmosphereSample) ([]byte, error) {
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	if err := w.Write([]string{"sample_id", "kiln_id", "observed_at", "temperature", "pressure"}); err != nil {
		return nil, err
	}
	for _, value := range values {
		if err := w.Write([]string{value.ID, value.KilnID, value.ObservedAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"), strconv.FormatFloat(value.Temperature, 'f', 3, 64), strconv.FormatFloat(value.Pressure, 'f', 3, 64)}); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return b.Bytes(), w.Error()
}

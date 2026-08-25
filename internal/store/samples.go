package store

import (
	"context"
	"database/sql"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"time"
)

func (d *DB) InsertSample(ctx context.Context, sample model.AtmosphereSample) error {
	if err := kilnExists(ctx, d.SQL, sample.KilnID); err != nil {
		return err
	}
	_, err := d.SQL.ExecContext(ctx, `INSERT INTO atmosphere_samples (id, kiln_id, run_id, temperature, pressure, oxidation, observed_at, sequence, source) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, sample.ID, sample.KilnID, sample.RunID, sample.Temperature, sample.Pressure, sample.Oxidation, timeValue(sample.ObservedAt), sample.Sequence, sample.Source)
	return wrap("insert sample", err)
}

func (d *DB) InsertSamples(ctx context.Context, samples []model.AtmosphereSample) error {
	return withTx(context.Background(), d.SQL, func(tx *sql.Tx) error {
		for _, sample := range samples {
			if err := kilnExistsTx(ctx, tx, sample.KilnID); err != nil {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO atmosphere_samples (id, kiln_id, run_id, temperature, pressure, oxidation, observed_at, sequence, source) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, sample.ID, sample.KilnID, sample.RunID, sample.Temperature, sample.Pressure, sample.Oxidation, timeValue(sample.ObservedAt), sample.Sequence, sample.Source); err != nil {
				continue
			}
		}
		return nil
	})
}

func scanSample(s scanner) (model.AtmosphereSample, error) {
	var x model.AtmosphereSample
	var observed string
	err := s.Scan(&x.ID, &x.KilnID, &x.RunID, &x.Temperature, &x.Pressure, &x.Oxidation, &observed, &x.Sequence, &x.Source)
	if err != nil {
		return x, err
	}
	x.ObservedAt, err = scanTime(observed)
	return x, err
}

func (d *DB) GetSample(ctx context.Context, id string) (model.AtmosphereSample, error) {
	x, err := scanSample(d.SQL.QueryRowContext(ctx, `SELECT id, kiln_id, run_id, temperature, pressure, oxidation, observed_at, sequence, source FROM atmosphere_samples WHERE id = ?`, id))
	if err != nil {
		return x, wrap("get sample", notFound(err))
	}
	return x, nil
}

func (d *DB) LatestSample(ctx context.Context, kilnID string) (model.AtmosphereSample, error) {
	x, err := scanSample(d.SQL.QueryRowContext(ctx, `SELECT id, kiln_id, run_id, temperature, pressure, oxidation, observed_at, sequence, source FROM atmosphere_samples WHERE kiln_id = ? ORDER BY observed_at DESC, sequence DESC LIMIT 1`, kilnID))
	if err != nil {
		return x, wrap("latest sample", notFound(err))
	}
	return x, nil
}

func (d *DB) ListSamples(ctx context.Context, filter model.SampleFilter) ([]model.AtmosphereSample, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	rows, err := d.SQL.QueryContext(ctx, `SELECT id, kiln_id, run_id, temperature, pressure, oxidation, observed_at, sequence, source FROM atmosphere_samples WHERE kiln_id = ? AND run_id = ? AND observed_at >= ? AND observed_at < ? ORDER BY observed_at, sequence LIMIT ?`, filter.KilnID, filter.RunID, timeValue(filter.Start), timeValue(filter.End), limit)
	if err != nil {
		return nil, wrap("list samples", err)
	}
	defer rows.Close()
	result := make([]model.AtmosphereSample, 0, limit)
	for rows.Next() {
		x, scanErr := scanSample(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, x)
	}
	return result, rows.Err()
}

func (d *DB) CountSamples(ctx context.Context, runID string) (int, error) {
	var n int
	err := d.SQL.QueryRowContext(ctx, `SELECT COUNT(1) FROM atmosphere_samples WHERE run_id = ?`, runID).Scan(&n)
	return n, wrap("count samples", err)
}
func kilnExistsTx(ctx context.Context, tx *sql.Tx, id string) error {
	var n int
	if err := tx.QueryRowContext(ctx, `SELECT 1 FROM kilns WHERE id = ?`, id).Scan(&n); err != nil {
		return notFound(err)
	}
	return nil
}

var _ = time.Time{}

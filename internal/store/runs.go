package store

import (
	"context"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"time"
)

func (d *DB) CreateRun(ctx context.Context, run model.FiringRun) error {
	if err := kilnExists(ctx, d.SQL, run.KilnID); err != nil {
		return err
	}
	_, err := d.SQL.ExecContext(ctx, `INSERT INTO firing_runs (id, kiln_id, status, started_at, finished_at, summary, score) VALUES (?, ?, ?, ?, ?, ?, ?)`, run.ID, run.KilnID, run.Status, timeValue(run.StartedAt), timeValue(run.FinishedAt), run.Summary, run.Score)
	return wrap("create run", err)
}

func scanRun(s scanner) (model.FiringRun, error) {
	var x model.FiringRun
	var started, finished string
	err := s.Scan(&x.ID, &x.KilnID, &x.Status, &started, &finished, &x.Summary, &x.Score)
	if err != nil {
		return x, err
	}
	x.StartedAt, err = scanTime(started)
	if err != nil {
		return x, err
	}
	x.FinishedAt, err = scanTime(finished)
	return x, err
}

func (d *DB) GetRun(ctx context.Context, id string) (model.FiringRun, error) {
	x, err := scanRun(d.SQL.QueryRowContext(ctx, `SELECT id, kiln_id, status, started_at, finished_at, summary, score FROM firing_runs WHERE id = ?`, id))
	if err != nil {
		return x, wrap("get run", notFound(err))
	}
	return x, nil
}
func (d *DB) UpdateRun(ctx context.Context, run model.FiringRun) error {
	r, err := d.SQL.ExecContext(ctx, `UPDATE firing_runs SET status = ?, finished_at = ?, summary = ?, score = ? WHERE id = ?`, run.Status, timeValue(run.FinishedAt), run.Summary, run.Score, run.ID)
	if err != nil {
		return wrap("update run", err)
	}
	n, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}
func (d *DB) ListRuns(ctx context.Context, kilnID string) ([]model.FiringRun, error) {
	rows, err := d.SQL.QueryContext(ctx, `SELECT id, kiln_id, status, started_at, finished_at, summary, score FROM firing_runs WHERE kiln_id = ? ORDER BY started_at DESC`, kilnID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.FiringRun, 0)
	for rows.Next() {
		x, scanErr := scanRun(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, x)
	}
	return result, rows.Err()
}
func (d *DB) RunWindow(ctx context.Context, runID string) (time.Time, time.Time, error) {
	var start, end string
	err := d.SQL.QueryRowContext(ctx, `SELECT MIN(observed_at), MAX(observed_at) FROM atmosphere_samples WHERE run_id = ?`, runID).Scan(&start, &end)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	a, err := scanTime(start)
	if err != nil {
		return a, time.Time{}, err
	}
	b, err := scanTime(end)
	return a, b, err
}

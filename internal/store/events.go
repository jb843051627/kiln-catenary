package store

import (
	"context"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"time"
)

func (d *DB) CreateEvent(ctx context.Context, event model.SafetyEvent) error {
	if err := kilnExists(ctx, d.SQL, event.KilnID); err != nil {
		return err
	}
	_, err := d.SQL.ExecContext(ctx, `INSERT INTO safety_events (id, run_id, kiln_id, kind, severity, message, created_at, resolved) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, event.ID, event.RunID, event.KilnID, event.Kind, event.Severity, event.Message, timeValue(event.CreatedAt), boolValue(event.Resolved))
	return wrap("create safety event", err)
}
func scanEvent(s scanner) (model.SafetyEvent, error) {
	var x model.SafetyEvent
	var created string
	var resolved int
	err := s.Scan(&x.ID, &x.RunID, &x.KilnID, &x.Kind, &x.Severity, &x.Message, &created, &resolved)
	x.Resolved = scanBool(resolved)
	if err != nil {
		return x, err
	}
	x.CreatedAt, err = scanTime(created)
	return x, err
}
func (d *DB) GetEvent(ctx context.Context, id string) (model.SafetyEvent, error) {
	x, err := scanEvent(d.SQL.QueryRowContext(ctx, `SELECT id, run_id, kiln_id, kind, severity, message, created_at, resolved FROM safety_events WHERE id = ?`, id))
	if err != nil {
		return x, wrap("get event", notFound(err))
	}
	return x, nil
}
func (d *DB) ResolveEvent(ctx context.Context, id string) error {
	r, err := d.SQL.ExecContext(ctx, `UPDATE safety_events SET resolved = 1 WHERE id = ? AND resolved = 0`, id)
	if err != nil {
		return err
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
func (d *DB) ListEvents(ctx context.Context, kilnID string, activeOnly bool) ([]model.SafetyEvent, error) {
	query := `SELECT id, run_id, kiln_id, kind, severity, message, created_at, resolved FROM safety_events WHERE kiln_id = ? ORDER BY created_at DESC`
	if activeOnly {
		query = `SELECT id, run_id, kiln_id, kind, severity, message, created_at, resolved FROM safety_events WHERE kiln_id = ? AND resolved = 0 ORDER BY created_at DESC`
	}
	rows, err := d.SQL.QueryContext(ctx, query, kilnID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.SafetyEvent, 0)
	for rows.Next() {
		x, scanErr := scanEvent(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, x)
	}
	return result, rows.Err()
}
func (d *DB) CountActiveEvents(ctx context.Context, kilnID string) (int, error) {
	var n int
	err := d.SQL.QueryRowContext(ctx, `SELECT COUNT(1) FROM safety_events WHERE kiln_id = ? AND resolved = 0 AND severity IN (?, ?)`, kilnID, model.SeverityWarn, model.SeverityAlarm).Scan(&n)
	return n, err
}
func (d *DB) EventsSince(ctx context.Context, kilnID string, since time.Time) ([]model.SafetyEvent, error) {
	rows, err := d.SQL.QueryContext(ctx, `SELECT id, run_id, kiln_id, kind, severity, message, created_at, resolved FROM safety_events WHERE kiln_id = ? AND created_at >= ? ORDER BY created_at`, kilnID, timeValue(since))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.SafetyEvent, 0)
	for rows.Next() {
		x, scanErr := scanEvent(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, x)
	}
	return result, rows.Err()
}

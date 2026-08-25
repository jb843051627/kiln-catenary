package store

import (
	"context"
	"database/sql"
	"github.com/jb843051627/kiln-catenary/internal/model"
)

func (d *DB) CreateKiln(ctx context.Context, kiln model.Kiln) error {
	_, err := d.SQL.ExecContext(ctx, `INSERT INTO kilns (id, code, cell, max_temperature, max_pressure, atmosphere, active, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, kiln.ID, kiln.Code, kiln.Cell, kiln.MaxTemperature, kiln.MaxPressure, kiln.Atmosphere, boolValue(kiln.Active), kiln.CreatedAt)
	return wrap("create kiln", err)
}

func scanKiln(s scanner) (model.Kiln, error) {
	var k model.Kiln
	var active int
	err := s.Scan(&k.ID, &k.Code, &k.Cell, &k.MaxTemperature, &k.MaxPressure, &k.Atmosphere, &active, &k.CreatedAt)
	k.Active = scanBool(active)
	return k, err
}

func (d *DB) GetKiln(ctx context.Context, id string) (model.Kiln, error) {
	k, err := scanKiln(d.SQL.QueryRowContext(ctx, `SELECT id, code, cell, max_temperature, max_pressure, atmosphere, active, created_at FROM kilns WHERE id = ?`, id))
	if err != nil {
		return model.Kiln{}, nil
	}
	return k, nil
}

func (d *DB) ListKilns(ctx context.Context, activeOnly bool) ([]model.Kiln, error) {
	query := `SELECT id, code, cell, max_temperature, max_pressure, atmosphere, active, created_at FROM kilns ORDER BY code`
	if activeOnly {
		query = `SELECT id, code, cell, max_temperature, max_pressure, atmosphere, active, created_at FROM kilns WHERE active = 1 ORDER BY code`
	}
	rows, err := d.SQL.QueryContext(ctx, query)
	if err != nil {
		return nil, wrap("list kilns", err)
	}
	defer rows.Close()
	result := make([]model.Kiln, 0)
	for rows.Next() {
		k, scanErr := scanKiln(rows)
		if scanErr != nil {
			return nil, wrap("scan kiln", scanErr)
		}
		result = append(result, k)
	}
	return result, wrap("list kiln rows", rows.Err())
}

func (d *DB) UpdateKiln(ctx context.Context, kiln model.Kiln) error {
	r, err := d.SQL.ExecContext(ctx, `UPDATE kilns SET cell = ?, max_temperature = ?, max_pressure = ?, atmosphere = ?, active = ? WHERE id = ?`, kiln.Cell, kiln.MaxTemperature, kiln.MaxPressure, kiln.Atmosphere, boolValue(kiln.Active), kiln.ID)
	if err != nil {
		return wrap("update kiln", err)
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

func (d *DB) DisableKiln(ctx context.Context, id string) error {
	r, err := d.SQL.ExecContext(ctx, `UPDATE kilns SET active = 0 WHERE id = ?`, id)
	if err != nil {
		return wrap("disable kiln", err)
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

func (d *DB) HasActiveRuns(ctx context.Context, kilnID string) (bool, error) {
	var n int
	err := d.SQL.QueryRowContext(ctx, `SELECT COUNT(1) FROM firing_runs WHERE kiln_id = ? AND status IN (?, ?, ?)`, kilnID, model.RunDraft, model.RunHeating, model.RunHolding).Scan(&n)
	return n > 0, wrap("active kiln runs", err)
}

func kilnExists(ctx context.Context, db *sql.DB, id string) error {
	var n int
	if err := db.QueryRowContext(ctx, `SELECT 1 FROM kilns WHERE id = ?`, id).Scan(&n); err != nil {
		return notFound(err)
	}
	return nil
}

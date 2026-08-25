package store

import (
	"context"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"time"
)

func (d *DB) CreateStage(ctx context.Context, stage model.ThermalStage) error {
	if err := kilnExists(ctx, d.SQL, stage.KilnID); err != nil {
		return err
	}
	_, err := d.SQL.ExecContext(ctx, `INSERT INTO thermal_stages (id, kiln_id, name, kind, sequence, start_temp, end_temp, hold_ns, status, interlock) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, stage.ID, stage.KilnID, stage.Name, stage.Kind, stage.Sequence, stage.StartTemp, stage.EndTemp, stage.Hold.Nanoseconds(), stage.Status, stage.Interlock)
	return wrap("create stage", err)
}

func scanStage(s scanner) (model.ThermalStage, error) {
	var x model.ThermalStage
	var hold int64
	err := s.Scan(&x.ID, &x.KilnID, &x.Name, &x.Kind, &x.Sequence, &x.StartTemp, &x.EndTemp, &hold, &x.Status, &x.Interlock)
	x.Hold = time.Duration(hold)
	return x, err
}

func (d *DB) GetStage(ctx context.Context, id string) (model.ThermalStage, error) {
	x, err := scanStage(d.SQL.QueryRowContext(ctx, `SELECT id, kiln_id, name, kind, sequence, start_temp, end_temp, hold_ns, status, interlock FROM thermal_stages WHERE id = ?`, id))
	if err != nil {
		return model.ThermalStage{}, wrap("get stage", err)
	}
	return x, nil
}

func (d *DB) ListStages(ctx context.Context, kilnID string) ([]model.ThermalStage, error) {
	rows, err := d.SQL.QueryContext(ctx, `SELECT id, kiln_id, name, kind, sequence, start_temp, end_temp, hold_ns, status, interlock FROM thermal_stages WHERE kiln_id = ? ORDER BY sequence`, kilnID)
	if err != nil {
		return nil, wrap("list stages", err)
	}
	defer rows.Close()
	result := make([]model.ThermalStage, 0)
	for rows.Next() {
		x, scanErr := scanStage(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, x)
	}
	return result, rows.Err()
}

func (d *DB) UpdateStageStatus(ctx context.Context, id, status string) error {
	r, err := d.SQL.ExecContext(ctx, `UPDATE thermal_stages SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return wrap("update stage", err)
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

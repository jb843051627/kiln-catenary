package store

import (
	"context"
	"database/sql"
	"github.com/jb843051627/kiln-catenary/internal/model"
)

func (d *DB) CreateZone(ctx context.Context, zone model.Zone) error {
	if err := kilnExists(ctx, d.SQL, zone.KilnID); err != nil {
		return wrap("zone kiln", err)
	}
	_, err := d.SQL.ExecContext(ctx, `INSERT INTO zones (id, kiln_id, name, target, ramp_rate, deviation_limit, enabled) VALUES (?, ?, ?, ?, ?, ?, ?)`, zone.ID, zone.KilnID, zone.Name, zone.Target, zone.RampRate, zone.DeviationLimit, boolValue(zone.Enabled))
	return wrap("create zone", err)
}

func scanZone(s scanner) (model.Zone, error) {
	var z model.Zone
	var enabled int
	err := s.Scan(&z.ID, &z.KilnID, &z.Name, &z.Target, &z.RampRate, &z.DeviationLimit, &enabled)
	z.Enabled = scanBool(enabled)
	return z, err
}

func (d *DB) GetZone(ctx context.Context, id string) (model.Zone, error) {
	z, err := scanZone(d.SQL.QueryRowContext(ctx, `SELECT id, kiln_id, name, target, ramp_rate, deviation_limit, enabled FROM zones WHERE id = ?`, id))
	if err != nil {
		return model.Zone{}, wrap("get zone", notFound(err))
	}
	return z, nil
}

func (d *DB) ListZones(ctx context.Context, kilnID string) ([]model.Zone, error) {
	rows, err := d.SQL.QueryContext(ctx, `SELECT id, kiln_id, name, target, ramp_rate, deviation_limit, enabled FROM zones WHERE kiln_id = ? ORDER BY name`, kilnID)
	if err != nil {
		return nil, wrap("list zones", err)
	}
	defer rows.Close()
	result := make([]model.Zone, 0)
	for rows.Next() {
		z, scanErr := scanZone(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, z)
	}
	return result, rows.Err()
}

func (d *DB) UpdateZone(ctx context.Context, zone model.Zone) error {
	r, err := d.SQL.ExecContext(ctx, `UPDATE zones SET name = ?, target = ?, ramp_rate = ?, deviation_limit = ?, enabled = ? WHERE id = ?`, zone.Name, zone.Target, zone.RampRate, zone.DeviationLimit, boolValue(zone.Enabled), zone.ID)
	if err != nil {
		return wrap("update zone", err)
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

func zoneExistsTx(ctx context.Context, tx *sql.Tx, id string) error {
	var n int
	if err := tx.QueryRowContext(ctx, `SELECT 1 FROM zones WHERE id = ?`, id).Scan(&n); err != nil {
		return notFound(err)
	}
	return nil
}

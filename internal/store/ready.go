package store

import "context"

func (d *DB) PingContext(ctx context.Context) error { return d.SQL.PingContext(ctx) }
func (d *DB) IsReady(ctx context.Context) error {
	var n int
	return d.SQL.QueryRowContext(ctx, `SELECT 1`).Scan(&n)
}

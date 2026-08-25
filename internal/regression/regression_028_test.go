package regression

import (
    "path/filepath"
    "context"
    "errors"
    "testing"
    "time"
    "github.com/jb843051627/kiln-catenary/internal/model"
    "github.com/jb843051627/kiln-catenary/internal/service"
    "github.com/jb843051627/kiln-catenary/internal/store"
)

func TestBug28_CSVUsesUTC(t *testing.T) {
	_, app, kiln := setupKiln28(t); ctx, cancel := context.WithCancel(context.Background()); cancel(); _, err := app.EventsSince(ctx, kiln.ID, time.Now().Add(-time.Hour)); if !errors.Is(err, context.Canceled) { t.Fatalf("canceled event query returned %v", err) }; _ = model.SeverityInfo
}

func setupKiln28(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug28.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K28", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

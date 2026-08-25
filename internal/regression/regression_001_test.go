package regression

import (
    "context"
    "errors"
    "path/filepath"
    "testing"
    "time"
    "github.com/jb843051627/kiln-catenary/internal/model"
    "github.com/jb843051627/kiln-catenary/internal/service"
    "github.com/jb843051627/kiln-catenary/internal/store"
)

func TestBug01_BatchSamplesHonorCancellation(t *testing.T) {
	db, app, kiln := setupKiln01(t)
	run, err := app.StartRun(context.Background(), kiln.ID)
	if err != nil { t.Fatal(err) }
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = app.BatchRecordSamples(ctx, kiln.ID, run.ID, []model.AtmosphereSample{{Temperature: 900, Pressure: .4, Oxidation: 20, ObservedAt: time.Now().UTC()}})
	if !errors.Is(err, context.Canceled) { t.Fatalf("canceled batch returned %v", err) }
	count, err := db.CountSamples(context.Background(), run.ID)
	if err != nil { t.Fatal(err) }
	if count != 0 { t.Fatalf("canceled batch wrote %d samples", count) }
}

func setupKiln01(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug01.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K01", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

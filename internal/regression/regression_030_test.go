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

func TestBug30_ZeroSampleLimitRejected(t *testing.T) {
	_, app, kiln := setupKiln30(t); ctx, cancel := context.WithCancel(context.Background()); cancel(); _, err := app.ListSamples(ctx, model.SampleFilter{KilnID: kiln.ID, RunID: "run", Start: time.Now().Add(-time.Hour), End: time.Now().Add(time.Hour), Limit: 1}); if !errors.Is(err, context.Canceled) { t.Fatalf("canceled sample query returned %v", err) }
}

func setupKiln30(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug30.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K30", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

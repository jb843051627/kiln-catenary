package regression

import (
    "path/filepath"
    "context"
    "testing"
    "time"
    "github.com/jb843051627/kiln-catenary/internal/model"
    "github.com/jb843051627/kiln-catenary/internal/service"
    "github.com/jb843051627/kiln-catenary/internal/store"
)

func TestBug17_RecentSamplesOwnStorage(t *testing.T) {
	_, app, kiln := setupKiln17(t)
	run, err := app.StartRun(context.Background(), kiln.ID); if err != nil { t.Fatal(err) }
	if _, err := app.RecordSample(context.Background(), model.AtmosphereSample{KilnID: kiln.ID, RunID: run.ID, Temperature: 900, Pressure: .4, Oxidation: 20, ObservedAt: time.Now().UTC()}); err != nil { t.Fatal(err) }
	first, err := app.RecentSamples(context.Background(), kiln.ID); if err != nil { t.Fatal(err) }; first[0].Temperature = 1
	second, err := app.RecentSamples(context.Background(), kiln.ID); if err != nil { t.Fatal(err) }; if second[0].Temperature == 1 { t.Fatal("recent sample cache was mutated") }
}

func setupKiln17(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug17.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K17", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

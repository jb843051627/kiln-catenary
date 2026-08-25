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

func TestBug20_BatchSamplesAreAtomic(t *testing.T) {
	db, _, kiln := setupKiln20(t); values := []model.AtmosphereSample{{ID: "same-sample", KilnID: kiln.ID, RunID: "run", Temperature: 900, Pressure: .4, Oxidation: 20, ObservedAt: time.Now().UTC(), Sequence: 1}, {ID: "same-sample", KilnID: kiln.ID, RunID: "run", Temperature: 901, Pressure: .4, Oxidation: 20, ObservedAt: time.Now().Add(time.Minute).UTC(), Sequence: 2}}
	if err := db.InsertSamples(context.Background(), values); err == nil { t.Fatal("duplicate batch succeeded") }; count, err := db.CountSamples(context.Background(), "run"); if err != nil { t.Fatal(err) }; if count != 0 { t.Fatalf("failed batch left %d samples", count) }
}

func setupKiln20(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug20.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K20", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

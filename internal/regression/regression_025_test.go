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

func TestBug25_StageEndIsExclusive(t *testing.T) {
	_, app, kiln := setupKiln25(t); stage, err := app.CreateStage(context.Background(), model.ThermalStage{KilnID: kiln.ID, Name: "hold", Kind: model.StageHold, Sequence: 1, Hold: time.Minute, Status: model.StageReady}); if err != nil { t.Fatal(err) }; ctx, cancel := context.WithCancel(context.Background()); cancel(); _, err = app.StageAvailableAt(ctx, stage.ID, time.Now(), time.Now().Add(-time.Minute), time.Now().Add(time.Minute)); if !errors.Is(err, context.Canceled) { t.Fatalf("canceled stage lookup returned %v", err) }
}

func setupKiln25(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug25.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K25", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

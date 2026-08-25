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

func TestBug11_ZeroHoldStageRejected(t *testing.T) {
	_, app, kiln := setupKiln11(t); stage, err := app.CreateStage(context.Background(), model.ThermalStage{KilnID: kiln.ID, Name: "hold", Kind: model.StageHold, Sequence: 1, Hold: time.Minute, Status: model.StageReady}); if err != nil { t.Fatal(err) }; if err := app.StartStage(context.Background(), stage.ID); err != nil { t.Fatal(err) }; err = app.StartStage(context.Background(), stage.ID); if !errors.Is(err, model.ErrInvalidState) { t.Fatalf("already running stage started with %v", err) }
}

func setupKiln11(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug11.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K11", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

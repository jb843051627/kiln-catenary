package regression

import (
    "context"
    "errors"
    "path/filepath"
    "testing"
    "github.com/jb843051627/kiln-catenary/internal/model"
    "github.com/jb843051627/kiln-catenary/internal/service"
    "github.com/jb843051627/kiln-catenary/internal/store"
)

func TestBug07_MissingStageRemainsClassifiable(t *testing.T) {
	_, app, kiln := setupKiln07(t)
	_, err := app.GetStage(context.Background(), "stage-missing")
	if !errors.Is(err, model.ErrNotFound) { t.Fatalf("missing stage for %s returned %v", kiln.ID, err) }
}

func setupKiln07(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug07.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K07", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

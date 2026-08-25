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

func TestBug09_ArchiveHonorsRunState(t *testing.T) {
	_, app, kiln := setupKiln09(t)
	run, err := app.StartRun(context.Background(), kiln.ID); if err != nil { t.Fatal(err) }
	_, err = app.ArchiveRun(context.Background(), run.ID)
	if !errors.Is(err, model.ErrInvalidState) { t.Fatalf("heating run archived with %v", err) }
}

func setupKiln09(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug09.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K09", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

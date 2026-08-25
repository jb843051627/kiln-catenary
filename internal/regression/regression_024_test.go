package regression

import (
    "path/filepath"
    "context"
    "errors"
    "testing"
    "github.com/jb843051627/kiln-catenary/internal/model"
    "github.com/jb843051627/kiln-catenary/internal/service"
    "github.com/jb843051627/kiln-catenary/internal/store"
)

func TestBug24_ZoneEnvelopeIsConsistent(t *testing.T) {
	_, app, kiln := setupKiln24(t); ok, err := app.AdmitSample(context.Background(), kiln.ID, 900, .4); if err != nil { t.Fatal(err) }; if !ok { t.Fatal("active kiln rejected a safe sample") }; if err := app.DisableKiln(context.Background(), kiln.ID); err != nil { t.Fatal(err) }; _, err = app.AdmitSample(context.Background(), kiln.ID, 900, .4); if !errors.Is(err, model.ErrConflict) { t.Fatalf("disabled kiln admitted sample with %v", err) }
}

func setupKiln24(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug24.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K24", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

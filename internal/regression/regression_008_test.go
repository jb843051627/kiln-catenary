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

func TestBug08_MissingZoneRemainsClassifiable(t *testing.T) {
	_, app, _ := setupKiln08(t)
	_, err := app.GetZone(context.Background(), "zone-missing")
	if !errors.Is(err, model.ErrNotFound) { t.Fatalf("missing zone returned %v", err) }
}

func setupKiln08(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug08.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K08", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

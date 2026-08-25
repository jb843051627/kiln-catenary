package regression

import (
    "testing"
    "path/filepath"
    "context"
    "errors"
    "github.com/jb843051627/kiln-catenary/internal/model"
    "github.com/jb843051627/kiln-catenary/internal/service"
    "github.com/jb843051627/kiln-catenary/internal/store"
)

func TestBug26_MissingKilnDoesNotBecomeEmpty(t *testing.T) {
	_, app, _ := setupKiln26(t); _, err := app.GetKiln(context.Background(), "kiln-missing"); if !errors.Is(err, model.ErrNotFound) { t.Fatalf("missing kiln became %v", err) }
}

func setupKiln26(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug26.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K26", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

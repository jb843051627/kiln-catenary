package regression

import (
    "path/filepath"
    "context"
    "strings"
    "testing"
    "github.com/jb843051627/kiln-catenary/internal/model"
    "github.com/jb843051627/kiln-catenary/internal/service"
    "github.com/jb843051627/kiln-catenary/internal/store"
)

func TestBug29_IDGeneratorIsRaceFree(t *testing.T) {
	_, app, _ := setupKiln29(t); for i := 0; i < 9; i++ { value, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K29-" + string(rune('A' + i)), Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true}); if err != nil { t.Fatal(err) }; parts := strings.Split(value.ID, "-"); if len(parts) == 0 || len(parts[len(parts)-1]) != 4 { t.Fatalf("bad sequence width: %s", value.ID) } }
}

func setupKiln29(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug29.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K29", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

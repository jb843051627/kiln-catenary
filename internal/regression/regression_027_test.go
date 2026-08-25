package regression

import (
    "testing"
    "path/filepath"
    "context"
    "errors"
    "time"
    "github.com/jb843051627/kiln-catenary/internal/model"
    "github.com/jb843051627/kiln-catenary/internal/service"
    "github.com/jb843051627/kiln-catenary/internal/store"
)

func TestBug27_WarningBlocksArchive(t *testing.T) {
	db, app, kiln := setupKiln27(t); run, err := app.StartRun(context.Background(), kiln.ID); if err != nil { t.Fatal(err) }; run.Status = model.RunEvaluated; run.FinishedAt = time.Now().UTC(); if err := db.UpdateRun(context.Background(), run); err != nil { t.Fatal(err) }; if err := db.CreateEvent(context.Background(), model.SafetyEvent{ID: "event-b27", RunID: run.ID, KilnID: kiln.ID, Kind: "interlock", Severity: model.SeverityWarn, Message: "watch", CreatedAt: time.Now().UTC()}); err != nil { t.Fatal(err) }; _, err = app.ArchiveRun(context.Background(), run.ID); if !errors.Is(err, model.ErrSafety) { t.Fatalf("warning event did not block archive: %v", err) }
}

func setupKiln27(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug27.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K27", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

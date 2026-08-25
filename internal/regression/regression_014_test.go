package regression

import (
    "context"
    "path/filepath"
    "testing"
    "time"
    "github.com/jb843051627/kiln-catenary/internal/model"
    "github.com/jb843051627/kiln-catenary/internal/service"
    "github.com/jb843051627/kiln-catenary/internal/store"
)

func TestBug14_SafetyEventKeepsAlarm(t *testing.T) {
	_, app, kiln := setupKiln14(t)
	run, err := app.StartRun(context.Background(), kiln.ID); if err != nil { t.Fatal(err) }
	if err := app.AdvanceRun(context.Background(), run.ID, model.RunHolding); err != nil { t.Fatal(err) }
	if err := app.AdvanceRun(context.Background(), run.ID, model.RunCooling); err != nil { t.Fatal(err) }
	_, err = app.RecordSample(context.Background(), model.AtmosphereSample{KilnID: kiln.ID, RunID: run.ID, Temperature: 1600, Pressure: .4, Oxidation: 20, ObservedAt: time.Now().UTC()}); if err != nil { t.Fatal(err) }
	if err := app.EvaluateRun(context.Background(), run.ID); err != nil { t.Fatal(err) }
	events, err := app.ListEvents(context.Background(), kiln.ID, true); if err != nil { t.Fatal(err) }
	if len(events) != 1 || events[0].Severity != model.SeverityAlarm { t.Fatalf("event severity = %#v", events) }
}

func setupKiln14(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug14.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K14", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

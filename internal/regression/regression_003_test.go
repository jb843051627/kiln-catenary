package regression

import (
    "context"
    "path/filepath"
    "sync"
    "testing"
    "time"
    "github.com/jb843051627/kiln-catenary/internal/model"
    "github.com/jb843051627/kiln-catenary/internal/service"
    "github.com/jb843051627/kiln-catenary/internal/store"
)

func TestBug03_SampleCacheIsRaceFree(t *testing.T) {
	_, app, kiln := setupKiln03(t)
	var wg sync.WaitGroup
	for worker := 0; worker < 20; worker++ {
		wg.Add(1)
		go func(offset int) { defer wg.Done(); for i := 0; i < 100; i++ { app.RefreshSampleCache(model.AtmosphereSample{KilnID: kiln.ID, RunID: "run-cache", Temperature: float64(800 + offset), Pressure: .4, Oxidation: 20, ObservedAt: time.Now().UTC(), Sequence: int64(i + 1)}); _, _ = app.LatestSample(context.Background(), kiln.ID) } }(worker)
	}
	wg.Wait()
}

func setupKiln03(t *testing.T) (*store.DB, *service.App, model.Kiln) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bug03.db"))
	if err != nil { t.Fatal(err) }
	app := service.NewApp(db)
	t.Cleanup(func() { app.Close(); db.Close() })
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "K03", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil { t.Fatal(err) }
	return db, app, kiln
}

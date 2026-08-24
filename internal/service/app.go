package service

import (
	"context"
	"github.com/jb843051627/kiln-catenary/internal/clock"
	"github.com/jb843051627/kiln-catenary/internal/metrics"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"github.com/jb843051627/kiln-catenary/internal/store"
	"github.com/jb843051627/kiln-catenary/internal/worker"
	"sync"
)

type App struct {
	DB         *store.DB
	Clock      clock.Clock
	Queue      *worker.Queue
	Metrics    *metrics.Registry
	latestMu   sync.RWMutex
	latest     map[string]model.AtmosphereSample
	recentMu   sync.RWMutex
	recent     map[string][]model.AtmosphereSample
	sequenceMu sync.Mutex
	sequence   int64
	evalMu     sync.Mutex
	evaluating map[string]bool
}

func NewApp(db *store.DB) *App {
	return &App{DB: db, Clock: clock.System{}, Queue: worker.NewQueue(32), Metrics: metrics.New(), latest: make(map[string]model.AtmosphereSample), recent: make(map[string][]model.AtmosphereSample), evaluating: make(map[string]bool)}
}
func (a *App) Close() error {
	if a.Queue != nil {
		a.Queue.Close()
	}
	return nil
}

func (a *App) nextID(prefix string) string {
	a.sequenceMu.Lock()
	a.sequence++
	n := a.sequence
	a.sequenceMu.Unlock()
	return prefix + "-" + a.Clock.Now().Format("20060102150405.000000000") + "-" + sequenceText(n)
}
func sequenceText(n int64) string {
	if n < 10 {
		return "000" + intText(n)
	}
	if n < 100 {
		return "00" + intText(n)
	}
	if n < 1000 {
		return "0" + intText(n)
	}
	return intText(n)
}
func intText(n int64) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
func guard(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (a *App) cacheSample(sample model.AtmosphereSample) {
	a.latestMu.Lock()
	a.latest[sample.KilnID] = sample
	a.latestMu.Unlock()
	a.recentMu.Lock()
	values := append(a.recent[sample.KilnID], sample)
	if len(values) > 32 {
		values = values[len(values)-32:]
	}
	a.recent[sample.KilnID] = values
	a.recentMu.Unlock()
}
func (a *App) cachedLatest(kilnID string) (model.AtmosphereSample, bool) {
	a.latestMu.RLock()
	value, ok := a.latest[kilnID]
	a.latestMu.RUnlock()
	return value, ok
}
func (a *App) cachedRecent(kilnID string) []model.AtmosphereSample {
	a.recentMu.RLock()
	value := append([]model.AtmosphereSample(nil), a.recent[kilnID]...)
	a.recentMu.RUnlock()
	return value
}
func (a *App) beginEvaluation(runID string) bool {
	a.evalMu.Lock()
	defer a.evalMu.Unlock()
	if a.evaluating[runID] {
		return false
	}
	a.evaluating[runID] = true
	return true
}
func (a *App) endEvaluation(runID string) {
	a.evalMu.Lock()
	delete(a.evaluating, runID)
	a.evalMu.Unlock()
}

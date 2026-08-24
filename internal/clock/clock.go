package clock

import (
	"sync"
	"time"
)

type Clock interface{ Now() time.Time }

type System struct{}

func (System) Now() time.Time { return time.Now().UTC() }

type Fixed struct {
	mu  sync.RWMutex
	now time.Time
}

func NewFixed(now time.Time) *Fixed { return &Fixed{now: now} }

func (f *Fixed) Now() time.Time { f.mu.RLock(); defer f.mu.RUnlock(); return f.now }

func (f *Fixed) Set(now time.Time) { f.mu.Lock(); f.now = now; f.mu.Unlock() }

func (f *Fixed) Advance(d time.Duration) { f.mu.Lock(); f.now = f.now.Add(d); f.mu.Unlock() }

package metrics

import "sync/atomic"

type Registry struct {
	samples atomic.Int64
	runs    atomic.Int64
	events  atomic.Int64
	errors  atomic.Int64
}

func New() *Registry        { return &Registry{} }
func (r *Registry) Sample() { r.samples.Add(1) }
func (r *Registry) Run()    { r.runs.Add(1) }
func (r *Registry) Event()  { r.events.Add(1) }
func (r *Registry) Error()  { r.errors.Add(1) }

type Snapshot struct {
	Samples int64 `json:"samples"`
	Runs    int64 `json:"runs"`
	Events  int64 `json:"events"`
	Errors  int64 `json:"errors"`
}

func (r *Registry) Snapshot() Snapshot {
	return Snapshot{Samples: r.samples.Load(), Runs: r.runs.Load(), Events: r.events.Load(), Errors: r.errors.Load()}
}

package worker

import (
	"context"
	"sync"
)

type Group struct {
	queue   *Queue
	mu      sync.Mutex
	pending int
}

func NewGroup(queue *Queue) *Group { return &Group{queue: queue} }

func (g *Group) Submit(ctx context.Context, job Job) error {
	g.mu.Lock()
	g.pending++
	g.mu.Unlock()
	return g.queue.Submit(ctx, func(jobCtx context.Context) error {
		defer func() { g.mu.Lock(); g.pending--; g.mu.Unlock() }()
		return job(jobCtx)
	})
}

func (g *Group) Pending() int { g.mu.Lock(); defer g.mu.Unlock(); return g.pending }

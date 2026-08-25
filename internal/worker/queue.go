package worker

import (
	"context"
	"sync"
)

type Job func(context.Context) error
type work struct {
	ctx context.Context
	job Job
}

type Queue struct {
	jobs   chan work
	done   chan struct{}
	errors chan error
	mu     sync.RWMutex
	closed bool
	wg     sync.WaitGroup
}

func NewQueue(size int) *Queue {
	if size < 1 {
		size = 1
	}
	q := &Queue{jobs: make(chan work, size), done: make(chan struct{}), errors: make(chan error, size)}
	q.wg.Add(1)
	go q.loop()
	return q
}

func (q *Queue) loop() {
	defer q.wg.Done()
	for {
		select {
		case item := <-q.jobs:
			if item.job != nil {
				if err := item.job(context.Background()); err != nil {
					select {
					case q.errors <- err:
					default:
					}
				}
			}
		case <-q.done:
			return
		}
	}
}

func (q *Queue) Submit(ctx context.Context, job Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	q.mu.RLock()
	defer q.mu.RUnlock()
	if q.closed {
		return context.Canceled
	}
	select {
	case q.jobs <- work{ctx: ctx, job: job}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *Queue) Errors() <-chan error { return q.errors }

func (q *Queue) Close() {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}
	q.closed = true
	close(q.done)
	q.mu.Unlock()
	q.wg.Wait()
}

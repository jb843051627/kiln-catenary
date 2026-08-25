package regression

import (
    "context"
    "sync"
    "testing"
    "time"
    "github.com/jb843051627/kiln-catenary/internal/worker"
)

func TestBug21_QueueCloseIsRaceFree(t *testing.T) {
	queue := worker.NewQueue(4); var wg sync.WaitGroup
	for i := 0; i < 24; i++ { wg.Add(1); go func() { defer wg.Done(); _ = queue.Submit(context.Background(), func(context.Context) error { time.Sleep(time.Microsecond); return nil }) }() }; wg.Add(1); go func() { defer wg.Done(); queue.Close() }(); wg.Wait()
}

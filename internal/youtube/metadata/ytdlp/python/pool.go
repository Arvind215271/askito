package python

import (
	"context"
	"sync"
	"time"
)

type WorkerExecution struct {
	VideoID   string        `json:"video_id"`
	StartedAt time.Time     `json:"started_at"`
	FinishedAt time.Time     `json:"finished_at"`
	Duration  time.Duration `json:"duration"`
}

type SinglePool struct {
	workers []*SingleClient
	taskCh  chan string
	wg      sync.WaitGroup
	once    sync.Once
}

func NewSinglePool(workerCount int) (*SinglePool, error) {
	workers := make([]*SingleClient, workerCount)
	for i := 0; i < workerCount; i++ {
		w, err := NewSingleWorker(i)
		if err != nil {
			return nil, err
		}
		workers[i] = w
	}
	return &SinglePool{workers: workers, taskCh: make(chan string)}, nil
}

func (p *SinglePool) Start(ctx context.Context, handler func(id string, data map[string]any, execution WorkerExecution, err error)) {
	for _, w := range p.workers {
		p.wg.Add(1)
		go func(worker *SingleClient) {
			defer p.wg.Done()
			for id := range p.taskCh {
				start := time.Now()
				data, err := worker.GetVideo(ctx, id)
				end := time.Now()
				execution := WorkerExecution{
					VideoID:   id,
					StartedAt: start,
					FinishedAt: end,
					Duration:  end.Sub(start),
				}
				handler(id, data, execution, err)
			}
		}(w)
	}
}

func (p *SinglePool) WarmUp(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(p.workers))
	for _, w := range p.workers {
		wg.Add(1)
		go func(worker *SingleClient) {
			defer wg.Done()
			if err := worker.WarmUp(ctx); err != nil {
				errCh <- err
			}
		}(w)
	}
	wg.Wait()
	close(errCh)
	if len(errCh) > 0 {
		return <-errCh
	}
	return nil
}

func (p *SinglePool) Submit(id string) {
	p.taskCh <- id
}

func (p *SinglePool) Close() {
	p.once.Do(func() {
		close(p.taskCh)
		p.wg.Wait()
		for _, w := range p.workers {
			w.Close()
		}
	})
}

type BatchPool struct {
	workers []*BatchClient
	taskCh  chan []string
	wg      sync.WaitGroup
	once    sync.Once
}

func NewBatchPool(workerCount int) (*BatchPool, error) {
	workers := make([]*BatchClient, workerCount)
	for i := 0; i < workerCount; i++ {
		w, err := NewBatchWorker(i)
		if err != nil {
			return nil, err
		}
		workers[i] = w
	}
	return &BatchPool{workers: workers, taskCh: make(chan []string)}, nil
}

func (p *BatchPool) Start(ctx context.Context, handler func(ids []string, results []map[string]any, duration time.Duration, err error)) {
	for _, w := range p.workers {
		p.wg.Add(1)
		go func(worker *BatchClient) {
			defer p.wg.Done()
			for batch := range p.taskCh {
				start := time.Now()
				results, err := worker.GetVideos(ctx, batch)
				duration := time.Since(start)
				handler(batch, results, duration, err)
			}
		}(w)
	}
}

func (p *BatchPool) Submit(batch []string) {
	p.taskCh <- batch
}

func (p *BatchPool) Close() {
	p.once.Do(func() {
		close(p.taskCh)
		p.wg.Wait()
		for _, w := range p.workers {
			w.Close()
		}
	})
}

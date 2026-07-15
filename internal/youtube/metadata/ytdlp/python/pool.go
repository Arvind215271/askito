package python

import (
	"context"
	"sync"
	"time"

	"github.com/Arvind215271/askito/internal/logger"
)

type WorkerExecution struct {
	VideoID    string        `json:"video_id"`
	StartedAt  time.Time     `json:"started_at"`
	FinishedAt time.Time     `json:"finished_at"`
	Duration   time.Duration `json:"duration"`
}

type SinglePool struct {
	manager *WorkerManager
	once    sync.Once
}

func NewSinglePool(
	workerCount int,
	log *logger.Logger,
) (*SinglePool, error) {
	manager, err := NewWorkerManager(
		workerCount,
		log,
	)
	if err != nil {
		return nil, err
	}

	return &SinglePool{
		manager: manager,
	}, nil
}

func (p *SinglePool) WarmUp(ctx context.Context) error {
	return p.manager.WarmUp(ctx)
}

func (p *SinglePool) GetVideo(ctx context.Context, videoID string) (map[string]any, error) {
	return p.manager.GetVideo(ctx, videoID)
}

func (p *SinglePool) GetSubtitle(ctx context.Context, videoID, language, subType, format string) ([]byte, error) {
	return p.manager.GetSubtitle(ctx, videoID, language, subType, format)
}

func (p *SinglePool) Close() {
	p.once.Do(func() {
		_ = p.manager.Close()
	})
}

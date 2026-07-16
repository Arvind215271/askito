package python

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Arvind215271/askito/internal/cache"
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
	cache   *cache.Manager
	once    sync.Once
}

func NewSinglePool(
	workerCount int,
	log *logger.Logger,
	cache *cache.Manager,
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
		cache:   cache,
	}, nil
}

func (p *SinglePool) WarmUp(ctx context.Context) error {
	return p.manager.WarmUp(ctx)
}

func (p *SinglePool) GetVideo(ctx context.Context, videoID string) (map[string]any, error) {
	key := p.cache.VideoKey()
	data, err := p.cache.Get(videoID, key)
	if err == nil {
		var result map[string]any
		if err := json.Unmarshal(data, &result); err == nil {
			return result, nil
		}
	}

	result, err := p.manager.GetVideo(ctx, videoID)
	if err != nil {
		return nil, err
	}

	// serialize and cache
	jsonData, _ := json.Marshal(result)
	_ = p.cache.Save(videoID, key, jsonData)

	return result, nil
}

func (p *SinglePool) GetPlaylist(
	ctx context.Context,
	playlistID string,
) (map[string]any, error) {
	key := p.cache.PlaylistKey()
	data, err := p.cache.Get(playlistID, key)
	if err == nil {
		var result map[string]any
		if err := json.Unmarshal(data, &result); err == nil {
			return result, nil
		}
	}

	result, err := p.manager.GetPlaylist(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	jsonData, _ := json.Marshal(result)
	_ = p.cache.Save(playlistID, key, jsonData)

	return result, nil
}

func (p *SinglePool) GetSubtitle(
	ctx context.Context,
	videoID,
	language,
	subType,
	format string,
) ([]byte, error) {
	key := p.cache.SubtitleKey(subType, language, format)
	// We use the path derived from the SubtitlePath method.
	// We append language and format because yt-dlp appends them.
	outputPath := p.cache.GetPath(videoID, p.cache.SubtitlePath(subType))

	content, err := p.cache.Get(videoID, key)
	if err == nil {
		return content, nil
	}

	_, err = p.manager.GetSubtitle(
		ctx,
		videoID,
		language,
		subType,
		format,
		outputPath,
	)
	if err != nil {
		return nil, err
	}

	// Reading the file after the worker writes it
	// Retry loop to handle potential filesystem lag
	var data []byte
	// for i := 0; i < 5; i++ {
	data, err = p.cache.Get(videoID, key)
		// if err == nil {
			// break
		// }
		// time.Sleep(100 * time.Millisecond)
	// }
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *SinglePool) Close() {
	p.once.Do(func() {
		_ = p.manager.Close()
	})
}

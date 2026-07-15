package python

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Arvind215271/askito/internal/logger"
)

type SingleClient struct {
	worker *PythonWorker
	logger *logger.Logger
}

func NewSingleWorker(
	workerID int,
	log *logger.Logger,
) (*SingleClient, error) {
	_, file, _, _ := runtime.Caller(0)

	script := filepath.Join(
		filepath.Dir(file),
		"python_worker_single.py",
	)

	worker, err := NewWorker(
		script,
		workerID,
		log,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create single worker %d: %w",
			workerID,
			err,
		)
	}

	return &SingleClient{
		worker: worker,
		logger: log.With(
			"worker_id", workerID,
		),
	}, nil
}

func (c *SingleClient) GetVideo(
	ctx context.Context,
	id string,
) (map[string]any, error) {
	c.logger.Debug(
		"fetching video metadata",
		"video_id", id,
	)

	if err := c.worker.SendCommand(
		VideoRequest{VideoID: id},
	); err != nil {
		return nil, fmt.Errorf(
			"send metadata request for video %s: %w",
			id,
			err,
		)
	}

	respData, err := c.worker.ReceiveResponse()
	if err != nil {
		return nil, fmt.Errorf(
			"receive metadata response for video %s: %w",
			id,
			err,
		)
	}

	var resp WorkerResponse

	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, fmt.Errorf(
			"decode metadata response for video %s: %w",
			id,
			err,
		)
	}

	if !resp.Ok {
		c.logger.Warn(
			"metadata request failed",
			"video_id", id,
			"error", resp.Err,
		)

		return nil, fmt.Errorf(
			"worker metadata error for video %s: %s",
			id,
			resp.Err,
		)
	}

	c.logger.Debug(
		"video metadata fetched",
		"video_id", id,
	)

	return resp.Data, nil
}

func getYTDLPCacheDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get current working directory: %w", err)
	}

	cacheDir := filepath.Join(
		cwd,
		".cache",
		"ytdlp",
	)

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf(
			"create yt-dlp cache directory %s: %w",
			cacheDir,
			err,
		)
	}

	return cacheDir, nil
}

func (c *SingleClient) GetSubtitle(
	ctx context.Context,
	videoID,
	language,
	subType,
	format string,
) ([]byte, error) {
	c.logger.Debug(
		"fetching subtitle",
		"video_id", videoID,
		"language", language,
		"type", subType,
		"format", format,
	)

	cacheDir, err := getYTDLPCacheDir()
	if err != nil {
		return nil, fmt.Errorf(
			"resolve yt-dlp cache directory: %w",
			err,
		)
	}

	req := SubtitleRequest{
		Cmd:      "subtitle",
		VideoID:  videoID,
		Language: language,
		Type:     subType,
		Format:   format,
		CacheDir: cacheDir,
	}

	if err := c.worker.SendCommand(req); err != nil {
		return nil, fmt.Errorf(
			"send subtitle request for video %s: %w",
			videoID,
			err,
		)
	}

	respData, err := c.worker.ReceiveResponse()
	if err != nil {
		return nil, fmt.Errorf(
			"receive subtitle response for video %s: %w",
			videoID,
			err,
		)
	}

	var resp SubtitleResponse

	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, fmt.Errorf(
			"decode subtitle response for video %s: %w",
			videoID,
			err,
		)
	}

	if !resp.Ok {
		c.logger.Warn(
			"subtitle request failed",
			"video_id", videoID,
			"language", language,
			"type", subType,
			"format", format,
			"error", resp.Err,
		)

		return nil, fmt.Errorf(
			"worker subtitle error for video %s: %s",
			videoID,
			resp.Err,
		)
	}

	if resp.Filename == "" {
		return nil, fmt.Errorf(
			"worker returned empty subtitle filename for video %s",
			videoID,
		)
	}

	subtitlePath := filepath.Join(
		cacheDir,
		videoID,
		resp.Filename,
	)

	content, err := os.ReadFile(subtitlePath)
	if err != nil {
		c.logger.Error(
			"failed to read downloaded subtitle",
			"video_id", videoID,
			"path", subtitlePath,
			"error", err,
		)

		return nil, fmt.Errorf(
			"read subtitle file %s: %w",
			subtitlePath,
			err,
		)
	}

	c.logger.Debug(
		"subtitle fetched",
		"video_id", videoID,
		"language", language,
		"format", format,
		"bytes", len(content),
	)

	return content, nil
}

func (c *SingleClient) WarmUp(ctx context.Context) error {
	if err := c.worker.SendCommand(WarmupRequest{Cmd: "warmup"}); err != nil {
		return err
	}
	respData, err := c.worker.ReceiveResponse()
	if err != nil {
		return err
	}
	var resp WorkerResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("worker warmup error: %s", resp.Err)
	}
	return nil
}

func (c *SingleClient) Close() error {
	return c.worker.Close()
}

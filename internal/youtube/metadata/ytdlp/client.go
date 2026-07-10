package ytdlp

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/Arvind215271/askito/internal/logger"
)

type Client struct {
	cache  *Cache
	logger *logger.Logger
}

const (
	maxRetries = 3
	timeout    = 30 * time.Second
)

func NewClient(cfg CacheConfig, logger *logger.Logger) *Client {
	return &Client{
		cache:  NewCache(cfg, logger),
		logger: logger,
	}
}

// Cleanup runs the cache cleanup routine.
func (c *Client) Cleanup() error {
	return c.cache.Cleanup()
}

// ValidateYTDLP checks if ytdlp is present.
func (c *Client) ValidateYTDLP() error {
	_, err := exec.LookPath("yt-dlp")
	return err
}

// Fetch fetches the request user asked from ytdlp...
func (c *Client) Fetch(ctx context.Context, args ...string) ([]byte, error) {
	var lastErr error

	c.logger.Debug("running yt-dlp", "args", args)

	for i := 0; i < maxRetries; i++ {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)

		cmd := exec.CommandContext(ctxWithTimeout, "yt-dlp", args...)
		output, err := cmd.CombinedOutput()

		cancel()

		if err == nil {
			c.logger.Debug("yt-dlp command succeeded")
			return output, nil
		}

		c.logger.Warn("yt-dlp attempt failed", "attempt", i+1, "error", err)
		lastErr = fmt.Errorf("%w\n%s", err, output)
	}

	c.logger.Error("yt-dlp command failed after retries", "error", lastErr)
	return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

func (c *Client) GetVideo(ctx context.Context, videoID string) (YTOutput, error) {
	// Try cache
	if cachedData, err := c.cache.Get(videoID); err == nil {
		var meta YTOutput
		if err := json.Unmarshal(cachedData, &meta); err == nil {
			return meta, nil
		}
	}

	url := "https://www.youtube.com/watch?v=" + videoID

	c.logger.Debug("fetching video metadata from ytdlp", "videoID", videoID)

	// Fetch if cache miss
	output, err := c.Fetch(
		ctx,
		"--skip-download",
		"--dump-single-json",
		url,
	)

	if err != nil {
		return YTOutput{}, err
	}

	c.logger.Info("fetched video metadata from ytdlp", "videoID", videoID)

	// Save to cache
	_ = c.cache.Save(videoID, output)

	var meta YTOutput
	err = json.Unmarshal(output, &meta)
	return meta, err
}

func (c *Client) GetPlaylist(ctx context.Context, playlistID string) (YTPlaylistOutput, error) {
	url := "https://www.youtube.com/playlist?list=" + playlistID

	c.logger.Debug("fetching playlist from ytdlp", "playlistID", playlistID)

	output, err := c.Fetch(
		ctx,
		"-j",
		"--flat-playlist",
		url,
	)
	if err != nil {
		return YTPlaylistOutput{}, err
	}

	c.logger.Info("fetched playlist from ytdlp", "playlistID", playlistID)

	var meta YTPlaylistOutput
	err = json.Unmarshal(output, &meta)
	return meta, err
}

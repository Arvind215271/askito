package ytdlp

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type Client struct{}

const (
	maxRetries = 3
	timeout    = 10 * time.Second
)

func NewClient() *Client {
	return &Client{}
}

func (c *Client) ValidateYTDLP() error {
	_, err := exec.LookPath("yt-dlp")
	return err
}

func (c *Client) Fetch(ctx context.Context, args ...string) ([]byte, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)

		cmd := exec.CommandContext(ctxWithTimeout, "yt-dlp", args...)
		output, err := cmd.CombinedOutput()

		cancel()

		if err == nil {
			return output, nil
		}

		lastErr = fmt.Errorf("%w\n%s", err, output)
	}

	return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

func (c *Client) GetVideo(ctx context.Context, videoID string) (YTOutput, error) {
	url := "https://www.youtube.com/watch?v=" + videoID

	output, err := c.Fetch(
		ctx,
		"--skip-download",
		"--dump-single-json",
		url,
	)
	if err != nil {
		return YTOutput{}, err
	}

	var meta YTOutput
	err = json.Unmarshal(output, &meta)
	return meta, err
}

func (c *Client) GetPlaylist(ctx context.Context, playlistID string) (YTPlaylistOutput, error) {
	url := "https://www.youtube.com/playlist?list=" + playlistID

	output, err := c.Fetch(
		ctx,
		"-j",
		"--flat-playlist",
		url,
	)
	if err != nil {
		return YTPlaylistOutput{}, err
	}

	var meta YTPlaylistOutput
	err = json.Unmarshal(output, &meta)
	return meta, err
}
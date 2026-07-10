package ytdlp

import (
	"context"

	"github.com/Arvind215271/askito/internal/logger"
	"github.com/Arvind215271/askito/internal/youtube"
)

type Provider struct {
	client *Client
	logger *logger.Logger
}

func NewProvider(client *Client, logger *logger.Logger) *Provider {
	return &Provider{
		client: client,
		logger: logger,
	}
}

func (p *Provider) GetVideo(ctx context.Context, videoID string) (youtube.Video, error) {
	p.logger.Debug("getting video from ytdlp provider", "videoID", videoID)
	meta, err := p.client.GetVideo(ctx, videoID)
	if err != nil {
		p.logger.Error("failed to get video from ytdlp provider", "error", err, "videoID", videoID)
		return youtube.Video{}, err
	}
	return MapVideo(meta), nil
}

func (p *Provider) GetPlaylist(ctx context.Context, playlistID string) (youtube.Playlist, error) {
	p.logger.Debug("getting playlist from ytdlp provider", "playlistID", playlistID)
	meta, err := p.client.GetPlaylist(ctx, playlistID)
	if err != nil {
		p.logger.Error("failed to get playlist from ytdlp provider", "error", err, "playlistID", playlistID)
		return youtube.Playlist{}, err
	}
	return MapPlaylist(meta), nil
}

package ytdlp

import (
	"context"
	"time"

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

func (p *Provider) GetPlaylistMetadata(ctx context.Context, playlistID string) (youtube.Playlist, error) {
	p.logger.Debug("getting playlist from ytdlp provider", "playlistID", playlistID)
	meta, err := p.client.GetPlaylist(ctx, playlistID)
	if err != nil {
		p.logger.Error("failed to get playlist from ytdlp provider", "error", err, "playlistID", playlistID)
		return youtube.Playlist{}, err
	}
	return MapPlaylist(meta), nil
}

func (p *Provider) GetPlaylistItems(ctx context.Context, playlistID string) ([]youtube.PlaylistItem, error) {
	meta, err := p.client.GetPlaylist(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	result := make([]youtube.PlaylistItem, 0, len(meta.Entries))
	for i, entry := range meta.Entries {
		result = append(result, youtube.PlaylistItem{
			VideoID:  entry.ID,
			Position: i,
			AddedAt:  time.Time{},
		})
	}
	return result, nil
}

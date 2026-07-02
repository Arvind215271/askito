package ytdlp

import (
	"context"

	"github.com/Arvind215271/askito/internal/youtube"
)

type Provider struct {
	client *Client
}

func NewProvider(client *Client) *Provider {
	return &Provider{client: client}
}

func (p *Provider) GetVideo(ctx context.Context, videoID string) (youtube.Video, error) {
	meta, err := p.client.GetVideo(ctx, videoID)
	if err != nil {
		return youtube.Video{}, err
	}
	return MapVideo(meta), nil
}

func (p *Provider) GetPlaylist(ctx context.Context, playlistID string) (youtube.Playlist, error) {
	meta, err := p.client.GetPlaylist(ctx, playlistID)
	if err != nil {
		return youtube.Playlist{}, err
	}
	return MapPlaylist(meta), nil
}

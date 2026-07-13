// internal/youtube/providers/youtube_api/provider.go

package youtubeapi

import (
	"context"

	"github.com/Arvind215271/askito/internal/logger"
	youtube "github.com/Arvind215271/askito/internal/youtube"
)

type Provider struct {
	client *Client
	logger *logger.Logger
}

func NewProvider(
	client *Client,
	logger *logger.Logger,
) *Provider {
	return &Provider{
		client: client,
		logger: logger,
	}
}

func (p *Provider) GetPlaylistMetadata(
	ctx context.Context,
	playlistID string,
) (youtube.Playlist, error) {

	// contain playlist metadata only not any videoID in it.
	playlist, err := p.client.GetPlaylist(
		ctx,
		playlistID,
	)
	if err != nil {
		return youtube.Playlist{}, err
	}

	return youtube.Playlist{
		ID:           playlist.Id,
		Title:        playlist.Snippet.Title,
		Description:  playlist.Snippet.Description,
		ChannelID:    playlist.Snippet.ChannelId,
		ChannelTitle: playlist.Snippet.ChannelTitle,

		ThumbnailURL: p.getPlaylistThumbnail(
			playlist,
		),

		ItemCount: int(playlist.ContentDetails.ItemCount),

		PrivacyStatus: playlist.Status.PrivacyStatus,

		PublishedAt: p.parseTime(
			playlist.Snippet.PublishedAt,
		),
	}, nil
}

// for a single video metadata fetching
func (p *Provider) GetVideo(
	ctx context.Context,
	videoID string,
) (youtube.Video, error) {

	videoList, err := p.client.GetVideos(
		ctx,
		[]string{videoID},
	)
	if err != nil {
		return youtube.Video{}, err
	}

	if len(videoList) == 0 {
		return youtube.Video{}, youtube.Err.Video.NotFound()
	}

	video := videoList[0]

	return p.mapVideo(video), nil
}

func (p *Provider) GetPlaylistItems(
	ctx context.Context,
	playlistID string,
) ([]youtube.PlaylistItem, error) {
	items, err := p.client.GetPlaylistItems(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	result := make([]youtube.PlaylistItem, 0, len(items))
	for _, item := range items {
		if item == nil || item.Snippet == nil || item.ContentDetails == nil {
			continue
		}
		result = append(result, youtube.PlaylistItem{
			VideoID:  item.ContentDetails.VideoId,
			Position: int(item.Snippet.Position),
			AddedAt:  p.parseTime(item.Snippet.PublishedAt),
		})
	}
	return result, nil
}

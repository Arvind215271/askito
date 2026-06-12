// internal/youtube/providers/youtube_api/provider.go

package youtubeapi

import (
	"context"

	youtube "github.com/Arvind215271/askito/internal/youtube"
)

type Provider struct {
	client *Client
}

func NewProvider(
	client *Client,
) *Provider {
	return &Provider{
		client: client,
	}
}

func (p *Provider) GetPlaylist(
	ctx context.Context,
	playlistID string,
) (youtube.Playlist, error) {

	playlist, err := p.client.GetPlaylist(
		ctx,
		playlistID,
	)
	if err != nil {
		return youtube.Playlist{}, err
	}

	items, err := p.client.GetPlaylistItems(
		ctx,
		playlistID,
	)
	if err != nil {
		return youtube.Playlist{}, err
	}

	videoIDs := make([]string, 0, len(items))

	for _, item := range items {
		if item == nil || item.ContentDetails == nil {
			continue
		}

		videoIDs = append(
			videoIDs,
			item.ContentDetails.VideoId,
		)
	}

	videos, err := p.client.GetVideos(
		ctx,
		videoIDs,
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

		ThumbnailURL: getPlaylistThumbnail(
			playlist,
		),

		ItemCount: int(playlist.ContentDetails.ItemCount),


		PrivacyStatus: playlist.Status.PrivacyStatus,

		PublishedAt: parseTime(
			playlist.Snippet.PublishedAt,
		),

		Videos: MapPlaylistVideos(
			items,
			videos,
		),
	}, nil
}


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
	
	return MapVideo(video), nil
}
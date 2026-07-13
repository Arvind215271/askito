package metadata

import (
	"context"

	"github.com/Arvind215271/askito/internal/youtube"
)

type Provider interface {
	GetPlaylistMetadata(
		ctx context.Context,
		playlistID string,
	) (youtube.Playlist, error)

	GetPlaylistItems(
		ctx context.Context,
		playlistID string,
	) ([]youtube.PlaylistItem, error)

	GetVideo(
		ctx context.Context,
		videoID string,
	) (youtube.Video, error)
}

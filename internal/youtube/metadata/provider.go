package metadata

import (
	"context"

	"github.com/Arvind215271/askito/internal/youtube"
)

type Provider interface {
	GetPlaylist(
		ctx context.Context,
		playlistID string,
	) (youtube.Playlist, error)

	GetVideo(
		ctx context.Context,
		videoID string,
	) (youtube.Video, error)
}

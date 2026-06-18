package youtube

import (
	"context"
)

type Provider interface {
    GetPlaylist(
        ctx context.Context,
        playlistID string,
    ) (Playlist, error)

    GetVideo(
        ctx context.Context,
        videoID string,
    ) (Video, error)
}
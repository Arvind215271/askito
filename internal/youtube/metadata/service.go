package metadata

import (
	"context"

	"github.com/Arvind215271/askito/internal/youtube"
)

type Service struct {
	apiProvider   Provider
	ytdlpProvider Provider
}

func NewService(
	apiProvider Provider,
	ytdlpProvider Provider,
) *Service {
	return &Service{
		apiProvider:   apiProvider,
		ytdlpProvider: ytdlpProvider,
	}
}

func (s *Service) GetPlaylist(
	ctx context.Context,
	playlistID string,
	providerType ProviderType,
) (youtube.Playlist, error) {
	if providerType == ProviderAPI {
		return s.apiProvider.GetPlaylist(ctx, playlistID)
	}
	return s.ytdlpProvider.GetPlaylist(ctx, playlistID)
}

func (s *Service) GetVideo(
	ctx context.Context,
	videoID string,
	providerType ProviderType,
) (youtube.Video, error) {
	if providerType == ProviderAPI {
		return s.apiProvider.GetVideo(ctx, videoID)
	}
	return s.ytdlpProvider.GetVideo(ctx, videoID)
}

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

func (s *Service) GetPlaylistMetadata(
	ctx context.Context,
	playlistID string,
	providerType ProviderType,
) (youtube.Playlist, error) {
	if providerType == ProviderAPI {
		return s.apiProvider.GetPlaylistMetadata(ctx, playlistID)
	}
	return s.ytdlpProvider.GetPlaylistMetadata(ctx, playlistID)
}

func (s *Service) GetPlaylistItems(
	ctx context.Context,
	playlistID string,
	providerType ProviderType,
) ([]youtube.PlaylistItem, error) {
	if providerType == ProviderAPI {
		return s.apiProvider.GetPlaylistItems(ctx, playlistID)
	}
	return s.ytdlpProvider.GetPlaylistItems(ctx, playlistID)
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

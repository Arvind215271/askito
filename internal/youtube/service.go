// internal/youtube/service.go

package youtube

import(
	"context"
)

type Service struct {
	apiProvider    Provider
	ytdlpProvider  Provider
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
) (Playlist, error) {
	if providerType == ProviderAPI {
		return s.apiProvider.GetPlaylist(ctx, playlistID)
	}
	return s.ytdlpProvider.GetPlaylist(ctx, playlistID)
}

func (s *Service) GetVideo(
	ctx context.Context,
	videoID string,
	providerType ProviderType,
) (Video, error) {
	if providerType == ProviderAPI {
		return s.apiProvider.GetVideo(ctx, videoID)
	}
	return s.ytdlpProvider.GetVideo(ctx, videoID)
}

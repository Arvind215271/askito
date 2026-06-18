// internal/youtube/service.go

package youtube

import(
	"context"
)

type Service struct {
	provider Provider
}

func NewService(
	provider Provider,
) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) GetPlaylist(
	ctx context.Context,
	playlistID string,
) (Playlist, error) {

	return s.provider.GetPlaylist(
		ctx,
		playlistID,
	)
}

func (s *Service) GetVideo(
	ctx context.Context,
	videoID string,
) (Video, error) {

	return s.provider.GetVideo(
		ctx,
		videoID,
	)
}
package transcript

import (
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



func (s *Service) Get(
    ctx context.Context,
    videoID string,
) (*Transcript, error) {

    return s.provider.GetTranscript(
        ctx,
        videoID,
    )
}

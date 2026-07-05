package transcript

import (
	"fmt"

	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Parse(result *subtitle.SubtitleResult) (*Transcript, error) {
	if result == nil {
		return nil, fmt.Errorf("subtitle result is nil")
	}

	switch result.Format {

	case "json3":
		segments, err := ParseJSON3ToSegments(result.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse json3 transcript: %w", err)
		}

		return &Transcript{
			Language: result.Language,
			Segments: segments,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported subtitle format: %s", result.Format)
	}
}
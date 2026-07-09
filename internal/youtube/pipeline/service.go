package pipeline

import (
	"context"
	"fmt"

	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/description"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube/signal"
	wordstats "github.com/Arvind215271/askito/internal/youtube/signal/word_stats"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
)

type Service struct {
	metadataService    *metadata.Service
	descriptionService *description.Service
	subtitleService    *subtitle.SubtitleService
	transcriptService  *transcript.Service
	signalService      *signal.SignalService
}

func NewService(
	metadataService *metadata.Service,
	descriptionService *description.Service,
	subtitleService *subtitle.SubtitleService,
	transcriptService *transcript.Service,
	signalService *signal.SignalService,
) *Service {
	return &Service{
		metadataService:    metadataService,
		descriptionService: descriptionService,
		subtitleService:    subtitleService,
		transcriptService:  transcriptService,
		signalService:      signalService,
	}
}

func (s *Service) Process(ctx context.Context, videoID string, planner *Planner) (*youtube.Video, error) {
	// var err error
	// 1. Metadata (Base) - Always required as per instructions that everything depends on it
	meta, err := s.metadataService.GetVideo(ctx, videoID, metadata.ProviderYTDLP)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch base metadata: %w", err)
	}
	video := &meta

	// 2. Description
	if planner.NeedsDescription() {
		descMeta, err := s.descriptionService.GetDescription(ctx, video.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch description: %w", err)
		}
		MapDescription(video, &descMeta)
	}

	// 3. Transcript & Subtitle
	var trans *transcript.Transcript
	if planner.NeedsTranscript() || planner.NeedsSignal() {
		sub, err := s.subtitleService.DownloadSubtitle(ctx, subtitle.DownloadRequest{VideoID: video.ID, Language: "en", Type: "automatic", }, video.SubtitleMetadata)
		if err != nil {
			return nil, fmt.Errorf("subtitle error: %w", err)
		}

		trans, err = s.transcriptService.Parse(sub)
		if err != nil {
			return nil, fmt.Errorf("transcript error: %w", err)
		}
		MapTranscript(video, trans)
	}

	// 4. Signal
	if planner.NeedsSignal() && trans != nil {
		sig := s.signalService.AnalyzeWordStats(trans, wordstats.DefaultWordStatsConfig())
		MapSignal(video, &sig)
	}

	return video, nil
}

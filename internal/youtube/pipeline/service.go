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

func (s *Service) Process(
	ctx context.Context,
	videoID string,
	req *Request,
	planner *Planner,
) (*youtube.Video, error) {
	if req == nil {
		return nil, fmt.Errorf("pipeline request is nil")
	}

	// 1. Metadata
	meta, err := s.metadataService.GetVideo(
		ctx,
		videoID,
		metadata.ProviderYTDLP,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch base metadata: %w", err)
	}

	video := &meta

	// 2. Description
	if planner.NeedsDescription() {
		descMeta, err := s.descriptionService.GetDescription(
			ctx,
			video.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch description: %w", err)
		}

		MapDescription(video, &descMeta)
	}

	// 3. Subtitle + Transcript
	var trans *transcript.Transcript

	if planner.NeedsTranscript() || planner.NeedsSignal() {
		// Defaults belong to the pipeline.
		// Use the caller-provided request when available.
		subReq := req.Subtitle
		if subReq == nil {
			defaultReq := subtitle.DefaultDownloadRequest(video.ID)
			subReq = &defaultReq
		}

		if err := subReq.Validate(); err != nil {
			return nil, fmt.Errorf("invalid subtitle request: %w", err)
		}

		sub, err := s.subtitleService.DownloadSubtitle(
			ctx,
			*subReq,
			video.SubtitleMetadata,
		)


		if err != nil {
			return nil, fmt.Errorf("subtitle error: %w", err)
		}

		trans, err = s.transcriptService.Parse(sub)
		if err != nil {
			return nil, fmt.Errorf("transcript error: %w", err)
		}

		// Transcript processing is controlled by the request.
		if req.Transcript != nil {
			processed, err := s.transcriptService.Process(
				trans,
				req.Transcript,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"transcript processing error: %w",
					err,
				)
			}

			video.Transcript = trans
			video.TranscriptText = processed
		} else {
			// Pipeline default behavior.
			MapTranscript(video, trans)
		}
	}

	// 4. Signal
	if planner.NeedsSignal() {
		sigReq := req.Signal
		if sigReq == nil {
			defaultReq := signal.DefaultSignalRequest(video.ID)
			sigReq = &defaultReq
		}

		if err := sigReq.Validate(); err != nil {
			return nil, fmt.Errorf("invalid signal request: %w", err)
		}

		// The current signal implementation still uses the existing
		// word-stats configuration. The request should be wired into
		// the actual signal analysis implementation as that API evolves.
		sig := s.signalService.AnalyzeWordStats(
			trans,
			wordstats.DefaultWordStatsConfig(),
		)

		MapSignal(video, &sig)
	}

	return video, nil
}
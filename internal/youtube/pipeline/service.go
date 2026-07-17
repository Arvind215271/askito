package pipeline

import (
	"context"
	"fmt"
	"sync"

	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/description"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube/signal"
	wordstats "github.com/Arvind215271/askito/internal/youtube/signal/word_stats"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
)

const (
	VideoErrorMetadata = "metadata_fetch_failed"
	VideoErrorSubtitle = "subtitle_fetch_failed"
)

type Service struct {
	metadataService    *metadata.Service
	descriptionService *description.Service
	subtitleService    *subtitle.SubtitleService
	transcriptService  *transcript.Service
	signalService      *signal.SignalService

	concurrency int
}

func NewService(
	metadataService *metadata.Service,
	descriptionService *description.Service,
	subtitleService *subtitle.SubtitleService,
	transcriptService *transcript.Service,
	signalService *signal.SignalService,
	concurrency int,
) *Service {
	if concurrency <= 0 {
		concurrency = 1
	}

	return &Service{
		metadataService:    metadataService,
		descriptionService: descriptionService,
		subtitleService:    subtitleService,
		transcriptService:  transcriptService,
		signalService:      signalService,
		concurrency:        concurrency,
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

	meta, err := s.metadataService.GetVideo(
		ctx,
		videoID,
		metadata.ProviderYTDLP,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch base metadata: %w", err)
	}

	video := &meta

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

	var trans *transcript.Transcript

	if planner.NeedsTranscript() || planner.NeedsSignal() {
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
			MapTranscript(video, trans)
		}
	}

	if planner.NeedsSignal() {
		sigReq := req.Signal
		if sigReq == nil {
			defaultReq := signal.DefaultSignalRequest(video.ID)
			sigReq = &defaultReq
		}

		if err := sigReq.Validate(); err != nil {
			return nil, fmt.Errorf("invalid signal request: %w", err)
		}

		sig := s.signalService.AnalyzeWordStats(
			trans,
			wordstats.DefaultWordStatsConfig(),
		)

		MapSignal(video, &sig)
	}

	return video, nil
}


func (s *Service) ProcessFaultTolerant(
	ctx context.Context,
	videoID string,
	req *Request,
	planner *Planner,
) *youtube.Video {
	video := &youtube.Video{
		ID: videoID,
	}

	if req == nil {
		video.Errors = append(
			video.Errors,
			"metadata fetch failed: pipeline request is nil",
		)

		return video
	}

	// 1. Metadata
	meta, err := s.metadataService.GetVideo(
		ctx,
		videoID,
		metadata.ProviderYTDLP,
	)
	if err != nil {
		video.Errors = append(
			video.Errors,
			fmt.Sprintf("metadata fetch failed: %v", err),
		)
	} else {
		*video = meta
	}

	// 2. Description
	if planner.NeedsDescription() {
		descMeta, err := s.descriptionService.GetDescription(
			ctx,
			video.Description,
		)
		if err != nil {
			video.Errors = append(
				video.Errors,
				fmt.Sprintf("metadata fetch failed: %v", err),
			)
		} else {
			MapDescription(video, &descMeta)
		}
	}

	// 3. Subtitle + Transcript
	var trans *transcript.Transcript

	if planner.NeedsTranscript() || planner.NeedsSignal() {
		subReq := req.Subtitle
		if subReq == nil {
			defaultReq := subtitle.DefaultDownloadRequest(video.ID)
			subReq = &defaultReq
		}

		if err := subReq.Validate(); err != nil {
			video.Errors = append(
				video.Errors,
				fmt.Sprintf("subtitle fetch failed: %v", err),
			)
		} else {
			sub, err := s.subtitleService.DownloadSubtitle(
				ctx,
				*subReq,
				video.SubtitleMetadata,
			)
			if err != nil {
				video.Errors = append(
					video.Errors,
					fmt.Sprintf("subtitle fetch failed: %v", err),
				)
			} else {
				trans, err = s.transcriptService.Parse(sub)
				if err != nil {
					video.Errors = append(
						video.Errors,
						fmt.Sprintf("subtitle fetch failed: %v", err),
					)
				}
			}
		}

		// Transcript processing only happens if transcript parsing succeeded.
		if trans != nil {
			if req.Transcript != nil {
				processed, err := s.transcriptService.Process(
					trans,
					req.Transcript,
				)
				if err != nil {
					video.Errors = append(
						video.Errors,
						fmt.Sprintf("subtitle fetch failed: %v", err),
					)
				} else {
					video.Transcript = trans
					video.TranscriptText = processed
				}
			} else {
				MapTranscript(video, trans)
			}
		}
	}

	// 4. Signal
	if planner.NeedsSignal() && trans != nil {
		sig := s.signalService.AnalyzeWordStats(
			trans,
			wordstats.DefaultWordStatsConfig(),
		)

		MapSignal(video, &sig)
	}

	return video
}





func (s *Service) ProcessVideos(
	ctx context.Context,
	videoIDs []string,
	req *Request,
	planner *Planner,
) []*youtube.Video {
	if len(videoIDs) == 0 {
		return nil
	}

	results := make([]*youtube.Video, len(videoIDs))

	concurrency := s.concurrency
	if concurrency <= 0 {
		concurrency = 1
	}

	sem := make(chan struct{}, concurrency)

	var wg sync.WaitGroup

	for i, videoID := range videoIDs {
		wg.Add(1)

		go func(index int, id string) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				results[index] = &youtube.Video{
					ID: id,
					Errors: []string{
						"metadata fetch failed: context cancelled",
					},
				}
				return
			}
			defer func() {
				<-sem
			}()

			results[index] = s.ProcessFaultTolerant(
				ctx,
				id,
				req,
				planner,
			)
		}(i, videoID)
	}

	wg.Wait()

	return results
}
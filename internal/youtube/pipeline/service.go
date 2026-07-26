package pipeline

import (
	"context"
	"fmt"
	"sync"

	"github.com/Arvind215271/askito/internal/logger"
	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/description"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube/signal"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
)

type Service struct {
	metadataService    *metadata.Service
	descriptionService *description.Service
	subtitleService    *subtitle.SubtitleService
	transcriptService  *transcript.Service
	signalService      *signal.SignalService

	logger      *logger.Logger
	concurrency int
}

func NewService(
	metadataService *metadata.Service,
	descriptionService *description.Service,
	subtitleService *subtitle.SubtitleService,
	transcriptService *transcript.Service,
	signalService *signal.SignalService,
	logger *logger.Logger,
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
		logger:             logger,
		concurrency:        concurrency,
	}
}

func (s *Service) ProcessResource(
	ctx context.Context,
	rc *ResourceContext,
	faultTolerant bool,
) *youtube.Video {
	if rc.Video == nil {
		rc.Video = &youtube.Video{}
	}
	video := rc.Video

	if rc.Request == nil || rc.Request.ExecutionPlan == nil {
		err := fmt.Errorf("pipeline request or execution plan is nil")
		if faultTolerant {
			video.Errors = append(video.Errors, youtube.Error{Message: err.Error()})
			return video
		}
		return nil
	}

	plan := rc.Request.ExecutionPlan

	// 1. Metadata Stage
	if plan.NeedsMetadata() {
		err := ProcessMetadata(ctx, rc, s.metadataService)
		if err != nil {
			s.logError("metadata processing failed", "videoID", video.ID, "error", err)
			if faultTolerant {
				video.Errors = append(video.Errors, youtube.Error{Message: fmt.Sprintf("metadata fetch failed: %v", err)})
				return video
			}
			return nil
		}
	}

	// 2. Description Stage
	if plan.NeedsDescription() {
		err := ProcessDescription(ctx, rc, s.descriptionService)
		if err != nil {
			s.logError("description processing failed", "videoID", video.ID, "error", err)
			if faultTolerant {
				video.Errors = append(video.Errors, youtube.Error{Message: fmt.Sprintf("description processing failed: %v", err)})
			} else {
				return nil
			}
		}
	}

	// 3. Subtitle + Transcript Stage
	var trans *transcript.Transcript
	if plan.NeedsSubtitle() || plan.NeedsTranscript() || plan.NeedsSignal() {
		var err error
		trans, err = ProcessSubtitleAndTranscript(ctx, rc, s.subtitleService, s.transcriptService)
		if err != nil {
			s.logError("subtitle/transcript processing failed", "videoID", video.ID, "error", err)
			if faultTolerant {
				video.Errors = append(video.Errors, youtube.Error{Message: fmt.Sprintf("subtitle/transcript processing failed: %v", err)})
			} else {
				return nil
			}
		}
	}

	// 4. Signal Stage
	if plan.NeedsSignal() && trans != nil {
		err := ProcessSignal(ctx, rc, s.signalService, trans)
		if err != nil {
			s.logError("signal processing failed", "videoID", video.ID, "error", err)
			if faultTolerant {
				video.Errors = append(video.Errors, youtube.Error{Message: fmt.Sprintf("signal processing failed: %v", err)})
			} else {
				return nil
			}
		}
	}

	return video
}

func (s *Service) Process(
	ctx context.Context,
	videoID string,
	req *Request,
) (*youtube.Video, error) {
	if req == nil || req.ExecutionPlan == nil {
		return nil, fmt.Errorf("pipeline request or execution plan is nil for video %s", videoID)
	}

	rc := &ResourceContext{
		Video:   &youtube.Video{ID: videoID},
		Request: req,
		Plan:    req.ExecutionPlan,
	}

	video := s.ProcessResource(ctx, rc, false)
	if video == nil || len(video.Errors) > 0 {
		return nil, fmt.Errorf("failed to process video %s", videoID)
	}

	return video, nil
}

func (s *Service) ProcessFaultTolerant(
	ctx context.Context,
	videoID string,
	req *Request,
) *youtube.Video {
	if req == nil || req.ExecutionPlan == nil {
		return &youtube.Video{
			ID: videoID,
			Errors: []youtube.Error{{Message: "pipeline request or execution plan is nil"}},
		}
	}

	rc := &ResourceContext{
		Video:   &youtube.Video{ID: videoID},
		Request: req,
		Plan:    req.ExecutionPlan,
	}

	return s.ProcessResource(ctx, rc, true)
}

func (s *Service) ProcessVideos(
	ctx context.Context,
	videoIDs []string,
	req *Request,
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
					ID:     id,
					Errors: []youtube.Error{{Message: "context cancelled"}},
				}
				return
			}
			defer func() { <-sem }()

			results[index] = s.ProcessFaultTolerant(ctx, id, req)
		}(i, videoID)
	}

	wg.Wait()
	return results
}

func (s *Service) logDebug(msg string, args ...any) {
	if s.logger == nil {
		return
	}
	s.logger.Debug(msg, args...)
}

func (s *Service) logWarn(msg string, args ...any) {
	if s.logger == nil {
		return
	}
	s.logger.Warn(msg, args...)
}

func (s *Service) logError(msg string, args ...any) {
	if s.logger == nil {
		return
	}
	s.logger.Error(msg, args...)
}

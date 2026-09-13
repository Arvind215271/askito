package pipeline

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Arvind215271/askito/internal/logger"
	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/description"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube/signal"
	"github.com/Arvind215271/askito/internal/youtube/stats"
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

func cleanErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	s := err.Error()
	lines := strings.Split(s, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "File \"") ||
			strings.HasPrefix(trimmed, "Traceback") ||
			strings.HasPrefix(trimmed, "During handling") ||
			strings.HasPrefix(trimmed, "The above exception") ||
			(strings.HasPrefix(line, "    ") && !strings.Contains(trimmed, "Error")) {
			continue
		}
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) > 0 {
		return strings.Join(result, " — ")
	}
	return s
}

func (s *Service) ProcessResourceWithStats(
	ctx context.Context,
	rc *ResourceContext,
	faultTolerant bool,
) (*youtube.Video, stats.PipelineStats) {
	var pipelineStats stats.PipelineStats

	if rc.Video == nil {
		rc.Video = &youtube.Video{}
	}
	video := rc.Video

	if rc.Request == nil || rc.Request.ExecutionPlan == nil {
		err := fmt.Errorf("pipeline request or execution plan is nil")
		if faultTolerant {
			video.Errors = append(video.Errors, youtube.Error{Message: cleanErrorMessage(err)})
			return video, pipelineStats
		}
		return nil, pipelineStats
	}

	plan := rc.Request.ExecutionPlan

	// 1. Metadata Stage
	if plan.NeedsMetadata() {
		metaStats, err := ProcessMetadata(ctx, rc, s.metadataService)
		pipelineStats.Metadata.Add(metaStats)
		if err != nil {
			s.logError("metadata processing failed", "videoID", video.ID, "error", err)
			if faultTolerant {
				video.Errors = append(video.Errors, youtube.Error{Message: fmt.Sprintf("metadata fetch failed: %s", cleanErrorMessage(err))})
				return video, pipelineStats
			}
			return nil, pipelineStats
		}
	}

	// 2. Description Stage
	if plan.NeedsDescription() {
		err := ProcessDescription(ctx, rc, s.descriptionService)
		if err != nil {
			s.logError("description processing failed", "videoID", video.ID, "error", err)
			if faultTolerant {
				video.Errors = append(video.Errors, youtube.Error{Message: fmt.Sprintf("description processing failed: %s", cleanErrorMessage(err))})
			} else {
				return nil, pipelineStats
			}
		}
	}

	// 3. Subtitle + Transcript Stage
	var trans *transcript.Transcript
	if plan.NeedsSubtitle() || plan.NeedsTranscript() || plan.NeedsSignal() {
		var subStats stats.SubtitleStats
		var err error
		trans, subStats, err = ProcessSubtitleAndTranscript(ctx, rc, s.subtitleService, s.transcriptService)
		pipelineStats.Subtitle.Add(subStats)
		if err != nil {
			s.logError("subtitle/transcript processing failed", "videoID", video.ID, "error", err)
			if faultTolerant {
				video.Errors = append(video.Errors, youtube.Error{Message: fmt.Sprintf("subtitle/transcript processing failed: %s", cleanErrorMessage(err))})
			} else {
				return nil, pipelineStats
			}
		}
	}

	// 4. Signal Stage
	if plan.NeedsSignal() && trans != nil {
		err := ProcessSignal(ctx, rc, s.signalService, trans)
		if err != nil {
			s.logError("signal processing failed", "videoID", video.ID, "error", err)
			if faultTolerant {
				video.Errors = append(video.Errors, youtube.Error{Message: fmt.Sprintf("signal processing failed: %s", cleanErrorMessage(err))})
			} else {
				return nil, pipelineStats
			}
		}
	}

	return video, pipelineStats
}

func (s *Service) ProcessResource(
	ctx context.Context,
	rc *ResourceContext,
	faultTolerant bool,
) *youtube.Video {
	video, _ := s.ProcessResourceWithStats(ctx, rc, faultTolerant)
	return video
}

func (s *Service) ProcessWithStats(
	ctx context.Context,
	videoID string,
	req *Request,
) (*youtube.Video, stats.PipelineStats, error) {
	if req == nil || req.ExecutionPlan == nil {
		return nil, stats.PipelineStats{}, fmt.Errorf("pipeline request or execution plan is nil for video %s", videoID)
	}

	rc := &ResourceContext{
		Video:   &youtube.Video{ID: videoID},
		Request: req,
		Plan:    req.ExecutionPlan,
	}

	video, pipelineStats := s.ProcessResourceWithStats(ctx, rc, false)
	if video == nil || len(video.Errors) > 0 {
		return nil, pipelineStats, fmt.Errorf("failed to process video %s", videoID)
	}

	return video, pipelineStats, nil
}

func (s *Service) Process(
	ctx context.Context,
	videoID string,
	req *Request,
) (*youtube.Video, error) {
	video, _, err := s.ProcessWithStats(ctx, videoID, req)
	return video, err
}

func (s *Service) ProcessFaultTolerantWithStats(
	ctx context.Context,
	videoID string,
	req *Request,
) (*youtube.Video, stats.PipelineStats) {
	if req == nil || req.ExecutionPlan == nil {
		return &youtube.Video{
			ID:     videoID,
			Errors: []youtube.Error{{Message: "pipeline request or execution plan is nil"}},
		}, stats.PipelineStats{}
	}

	rc := &ResourceContext{
		Video:   &youtube.Video{ID: videoID},
		Request: req,
		Plan:    req.ExecutionPlan,
	}

	return s.ProcessResourceWithStats(ctx, rc, true)
}

func (s *Service) ProcessFaultTolerant(
	ctx context.Context,
	videoID string,
	req *Request,
) *youtube.Video {
	video, _ := s.ProcessFaultTolerantWithStats(ctx, videoID, req)
	return video
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

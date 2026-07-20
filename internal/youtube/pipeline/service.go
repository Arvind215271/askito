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

func (s *Service) Process(
	ctx context.Context,
	videoID string,
	req *Request,
	planner *Planner,
) (*youtube.Video, error) {

	if req == nil {
		return nil, fmt.Errorf(
			"pipeline request is nil for video %s",
			videoID,
		)
	}

	meta, err := s.metadataService.GetVideo(
		ctx,
		videoID,
		metadata.ProviderYTDLP,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to fetch base metadata for video %s: %w",
			videoID,
			err,
		)
	}

	video := &meta

	if planner.NeedsDescription() {
		descMeta, err := s.descriptionService.GetDescription(
			ctx,
			video.Description,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to fetch description for video %s: %w",
				videoID,
				err,
			)
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
			return nil, fmt.Errorf(
				"invalid subtitle request for video %s: %w",
				videoID,
				err,
			)
		}

		sub, err := s.subtitleService.DownloadSubtitle(
			ctx,
			*subReq,
			video.SubtitleMetadata,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"subtitle error for video %s: %w",
				videoID,
				err,
			)
		}

		trans, err = s.transcriptService.Parse(sub)
		if err != nil {
			return nil, fmt.Errorf(
				"transcript error for video %s: %w",
				videoID,
				err,
			)
		}

		if req.Transcript != nil {
			processed, err := s.transcriptService.Process(
				trans,
				req.Transcript,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"transcript processing error for video %s: %w",
					videoID,
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
			return nil, fmt.Errorf(
				"invalid signal request for video %s: %w",
				videoID,
				err,
			)
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
		err := fmt.Errorf(
			"pipeline request is nil",
		)

		s.logError(
			"pipeline request is nil",
			"videoID",
			videoID,
			"error",
			err,
		)

		video.Errors = append(
			video.Errors,
			fmt.Sprintf(
				"metadata fetch failed: %v",
				err,
			),
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

		s.logError(
			"metadata fetch failed",
			"videoID",
			videoID,
			"error",
			err,
		)

		video.Errors = append(
			video.Errors,
			fmt.Sprintf(
				"metadata fetch failed: %v",
				err,
			),
		)

	} else {

		*video = meta

		s.logDebug(
			"metadata fetched",
			"videoID",
			videoID,
		)
	}

	// 2. Description

	if planner.NeedsDescription() {

		s.logDebug(
			"fetching description metadata",
			"videoID",
			videoID,
		)

		descMeta, err := s.descriptionService.GetDescription(
			ctx,
			video.Description,
		)
		if err != nil {

			s.logError(
				"description processing failed",
				"videoID",
				videoID,
				"error",
				err,
			)

			video.Errors = append(
				video.Errors,
				fmt.Sprintf(
					"description processing failed: %v",
					err,
				),
			)

		} else {

			MapDescription(video, &descMeta)

			s.logDebug(
				"description processed",
				"videoID",
				videoID,
			)
		}
	}

	// 3. Subtitle + Transcript

	var trans *transcript.Transcript

	if planner.NeedsTranscript() || planner.NeedsSignal() {

		var subReq subtitle.DownloadRequest

		if req.Subtitle == nil {

			subReq = subtitle.DefaultDownloadRequest(
				video.ID,
			)

			s.logDebug(
				"using default subtitle request",
				"videoID",
				videoID,
				"language",
				subReq.Language,
				"type",
				subReq.Type,
				"format",
				subReq.Format,
			)

		} else {

			// Copy the shared request before setting the per-video ID.
			subReq = *req.Subtitle
			subReq.VideoID = video.ID

			s.logDebug(
				"using requested subtitle configuration",
				"videoID",
				videoID,
				"language",
				subReq.Language,
				"type",
				subReq.Type,
				"format",
				subReq.Format,
			)
		}

		if err := subReq.Validate(); err != nil {

			s.logError(
				"invalid subtitle request",
				"videoID",
				videoID,
				"error",
				err,
			)

			video.Errors = append(
				video.Errors,
				fmt.Sprintf(
					"subtitle request validation failed: %v",
					err,
				),
			)

		} else {

			s.logDebug(
				"downloading subtitle",
				"videoID",
				videoID,
				"language",
				subReq.Language,
				"type",
				subReq.Type,
				"format",
				subReq.Format,
			)

			sub, err := s.subtitleService.DownloadSubtitle(
				ctx,
				subReq,
				video.SubtitleMetadata,
			)

			if err != nil {

				s.logError(
					"subtitle fetch failed",
					"videoID",
					videoID,
					"language",
					subReq.Language,
					"type",
					subReq.Type,
					"format",
					subReq.Format,
					"error",
					err,
				)

				video.Errors = append(
					video.Errors,
					fmt.Sprintf(
						"subtitle fetch failed: %v",
						err,
					),
				)

			} else {

				s.logDebug(
					"subtitle fetched",
					"videoID",
					videoID,
					"bytes",
					len(sub.Content),
				)

				trans, err = s.transcriptService.Parse(sub)

				if err != nil {

					s.logError(
						"transcript parsing failed",
						"videoID",
						videoID,
						"error",
						err,
					)

					video.Errors = append(
						video.Errors,
						fmt.Sprintf(
							"transcript parsing failed: %v",
							err,
						),
					)

				} else {

					s.logDebug(
						"transcript parsed",
						"videoID",
						videoID,
					)
				}
			}
		}

		// Transcript processing

		if trans != nil {

			if req.Transcript != nil {

				s.logDebug(
					"processing transcript",
					"videoID",
					videoID,
				)

				processed, err := s.transcriptService.Process(
					trans,
					req.Transcript,
				)

				if err != nil {

					s.logError(
						"transcript processing failed",
						"videoID",
						videoID,
						"error",
						err,
					)

					video.Errors = append(
						video.Errors,
						fmt.Sprintf(
							"transcript processing failed: %v",
							err,
						),
					)

				} else {

					video.Transcript = trans
					video.TranscriptText = processed

					s.logDebug(
						"transcript processed",
						"videoID",
						videoID,
					)
				}

			} else {

				MapTranscript(video, trans)

				s.logDebug(
					"transcript mapped",
					"videoID",
					videoID,
				)
			}
		}
	}

	// 4. Signal

	if planner.NeedsSignal() && trans != nil {

		s.logDebug(
			"analyzing signal",
			"videoID",
			videoID,
		)

		sig := s.signalService.AnalyzeWordStats(
			trans,
			wordstats.DefaultWordStatsConfig(),
		)

		MapSignal(video, &sig)

		s.logDebug(
			"signal analysis completed",
			"videoID",
			videoID,
		)
	}

	if len(video.Errors) > 0 {

		s.logWarn(
			"video processing completed with errors",
			"videoID",
			videoID,
			"errorCount",
			len(video.Errors),
			"errors",
			video.Errors,
		)

	} else {

		s.logDebug(
			"video processing completed",
			"videoID",
			videoID,
		)
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

	results := make(
		[]*youtube.Video,
		len(videoIDs),
	)

	concurrency := s.concurrency

	if concurrency <= 0 {
		concurrency = 1
	}

	sem := make(
		chan struct{},
		concurrency,
	)

	var wg sync.WaitGroup

	for i, videoID := range videoIDs {

		wg.Add(1)

		go func(
			index int,
			id string,
		) {

			defer wg.Done()

			select {

			case sem <- struct{}{}:

				s.logDebug(
					"video processing started",
					"videoID",
					id,
					"index",
					index,
				)

			case <-ctx.Done():

				err := fmt.Errorf(
					"context cancelled",
				)

				s.logError(
					"video processing cancelled",
					"videoID",
					id,
					"error",
					err,
				)

				results[index] = &youtube.Video{
					ID: id,
					Errors: []string{
						fmt.Sprintf(
							"metadata fetch failed: %v",
							err,
						),
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

		}(
			i,
			videoID,
		)
	}

	wg.Wait()

	s.logDebug(
		"batch video processing completed",
		"videoCount",
		len(videoIDs),
		"concurrency",
		concurrency,
	)

	return results
}

func (s *Service) logDebug(
	msg string,
	args ...any,
) {
	if s.logger == nil {
		return
	}

	s.logger.Debug(msg, args...)
}

func (s *Service) logWarn(
	msg string,
	args ...any,
) {
	if s.logger == nil {
		return
	}

	s.logger.Warn(msg, args...)
}

func (s *Service) logError(
	msg string,
	args ...any,
) {
	if s.logger == nil {
		return
	}

	s.logger.Error(msg, args...)
}
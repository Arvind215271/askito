package video

import (
	"context"
	"net/http"

	"github.com/Arvind215271/askito/internal/api"
	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/description"
	youtubeurl "github.com/Arvind215271/askito/internal/youtube/input"
	"github.com/Arvind215271/askito/internal/youtube/signal"
	wordstats "github.com/Arvind215271/askito/internal/youtube/signal/word_stats"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	youtubeService    *youtube.Service
	subtitleService   *subtitle.SubtitleService
	transcriptService *transcript.Service
	signalService     *signal.SignalService
}

func NewHandler(youtubeService *youtube.Service, subtitleService *subtitle.SubtitleService, transcriptService *transcript.Service, signalService *signal.SignalService) *Handler {
	return &Handler{
		youtubeService:    youtubeService,
		subtitleService:   subtitleService,
		transcriptService: transcriptService,
		signalService:     signalService,
	}
}

func (h *Handler) GetSubtitleOptions(c *echo.Context) error {
	var req SubtitleOptionsRequest
	if err := (*c).Bind(&req); err != nil {
		return Err.BadRequest("Invalid request").Wrap(err)
	}

	if err := api.Validate(req); err != nil {
		return err
	}

	parsed, err := youtubeurl.Parse(req.URL)
	if err != nil || parsed.InputType != youtubeurl.InputTypeVideo {
		return Err.InvalidURL().Wrap(err)
	}

	video, err := h.getVideo((*c).Request().Context(), parsed.ID, youtube.ProviderYTDLP)
	if err != nil {
		return Err.FetchFailed(err)
	}

	return (*c).JSON(http.StatusOK, video.SubtitleMetadata)
}

func (h *Handler) DownloadSubtitle(c *echo.Context) error {
	var req SubtitleDownloadRequest
	if err := (*c).Bind(&req); err != nil {
		return Err.BadRequest("Invalid request").Wrap(err)
	}

	if err := api.Validate(req); err != nil {
		return err
	}

	parsed, err := youtubeurl.Parse(req.URL)
	if err == nil {
		req.URL = parsed.ID
	}

	video, err := h.getVideo((*c).Request().Context(), req.URL, youtube.ProviderYTDLP)
	if err != nil {
		return Err.FetchFailed(err)
	}

	result, err := h.subtitleService.DownloadSubtitle((*c).Request().Context(), subtitle.DownloadRequest{
		VideoID:  video.ID,
		Type:     req.Type,
		Language: req.Language,
		Format:   req.Format,
	}, video.SubtitleMetadata)
	if err != nil {
		return Err.InternalError(err)
	}

	(*c).Response().Header().Set("Content-Type", "application/octet-stream")
	(*c).Response().Header().Set("Content-Disposition", "attachment; filename=subtitle."+result.Format)
	return (*c).Blob(http.StatusOK, "application/octet-stream", result.Content)
}

func (h *Handler) GetTranscript(c *echo.Context) error {
	var req TranscriptRequest
	if err := (*c).Bind(&req); err != nil {
		return Err.BadRequest("Invalid request").Wrap(err)
	}

	if err := api.Validate(req); err != nil {
		return err
	}

	parsed, err := youtubeurl.Parse(req.URL)
	if err != nil {
		return Err.InvalidURL().Wrap(err)
	}

	video, err := h.getVideo((*c).Request().Context(), parsed.ID, youtube.ProviderYTDLP)
	if err != nil {
		return Err.FetchFailed(err)
	}

	result, err := h.subtitleService.DownloadSubtitle((*c).Request().Context(), subtitle.DownloadRequest{
		VideoID:  video.ID,
		Type:     req.Type,
		Language: req.Language,
		Format:   "json3",
	}, video.SubtitleMetadata)
	if err != nil {
		return Err.InternalError(err)
	}

	transcript, err := h.transcriptService.Parse(result)
	if err != nil {
		return Err.InternalError(err)
	}

	return (*c).JSON(http.StatusOK, transcript)
}

func (h *Handler) GetVideoByID(c *echo.Context) error {
	var req VideoByIDRequest
	if err := (*c).Bind(&req); err != nil {
		return Err.BadRequest("Invalid request").Wrap(err)
	}

	if err := api.Validate(req); err != nil {
		return err
	}

	providerType := youtube.ProviderType(req.Provider)

	video, err := h.getVideo((*c).Request().Context(), req.ID, providerType)
	if err != nil {
		return Err.FetchFailed(err)
	}

	return (*c).JSON(http.StatusOK, VideoResponse{Video: video})
}

func (h *Handler) GetVideoSignals(c *echo.Context) error {
	var req SignalRequest
	if err := (*c).Bind(&req); err != nil {
		return Err.BadRequest("Invalid request").Wrap(err)
	}

	if err := api.Validate(req); err != nil {
		return err
	}

	parsed, err := youtubeurl.Parse(req.URL)
	if err != nil {
		return Err.InvalidURL().Wrap(err)
	}

	video, err := h.getVideo((*c).Request().Context(), parsed.ID, youtube.ProviderYTDLP)
	if err != nil {
		return Err.FetchFailed(err)
	}

	result, err := h.subtitleService.DownloadSubtitle((*c).Request().Context(), subtitle.DownloadRequest{
		VideoID:  video.ID,
		Type:     req.Type,
		Language: req.Language,
		Format:   "json3",
	}, video.SubtitleMetadata)
	if err != nil {
		return Err.InternalError(err)
	}

	t, err := h.transcriptService.Parse(result)
	if err != nil {
		return Err.InternalError(err)
	}

	// Apply defaults
	windowSize := req.WindowSize
	if windowSize <= 0 {
		windowSize = 300
	}
	bucketCount := req.BucketCount
	if bucketCount <= 0 {
		if req.Analysis == "word-stats" {
			bucketCount = 32
		} else {
			bucketCount = 3
		}
	}

	cfg := wordstats.AnalysisConfig{
		UseHeavyStopWords: req.UseHeavy,
		MinFreq:           req.MinFreq,
		Depth:             req.Depth,
		WindowSize:        windowSize,
		BucketCount:       bucketCount,
	}

	resp := SignalResponse{URL: req.URL}
	switch req.Analysis {
	case "word-stats":
		stats := h.signalService.AnalyzeWordStats(t, cfg)
		resp.WordStats = &stats
	case "window-stats":
		stats := h.signalService.AnalyzeWindowedStats(t, cfg)
		resp.WindowStats = stats
	default:
		return Err.BadRequest("Invalid analysis type")
	}

	return (*c).JSON(http.StatusOK, resp)
}

func (h *Handler) GetVideoByURL(c *echo.Context) error {
	var req VideoRequest
	if err := (*c).Bind(&req); err != nil {
		return Err.BadRequest("Invalid request").Wrap(err)
	}

	if err := api.Validate(req); err != nil {
		return err
	}

	parsed, err := youtubeurl.Parse(req.URL)
	if err != nil {
		return Err.InvalidURL().Wrap(err)
	}

	if parsed.InputType != youtubeurl.InputTypeVideo {
		return Err.NotAVideo()
	}

	providerType := youtube.ProviderType(req.Provider)

	video, err := h.getVideo((*c).Request().Context(), parsed.ID, providerType)
	if err != nil {
		return Err.FetchFailed(err)
	}

	return (*c).JSON(http.StatusOK, VideoResponse{Video: video})
}

func (h *Handler) getVideo(ctx context.Context, id string, providerType youtube.ProviderType) (youtube.Video, error) {
	if providerType == "" {
		providerType = youtube.ProviderYTDLP
	}

	if providerType != youtube.ProviderAPI && providerType != youtube.ProviderYTDLP {
		return youtube.Video{}, Err.InvalidProvider()
	}

	video, err := h.youtubeService.GetVideo(ctx, id, providerType)
	if err != nil {
		return youtube.Video{}, err
	}

	video.DescriptionMetadata = description.ProcessDescription(video.Description)
	return video, nil
}

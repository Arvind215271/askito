package signal

import (
	"net/http"

	"github.com/Arvind215271/askito/internal/api"
	youtubeurl "github.com/Arvind215271/askito/internal/youtube/input"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube/signal"
	wordstats "github.com/Arvind215271/askito/internal/youtube/signal/word_stats"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	youtubeService    *metadata.Service
	subtitleService   *subtitle.SubtitleService
	transcriptService *transcript.Service
	signalService     *signal.SignalService
}

func NewHandler(youtubeService *metadata.Service, subtitleService *subtitle.SubtitleService, transcriptService *transcript.Service, signalService *signal.SignalService) *Handler {
	return &Handler{
		youtubeService:    youtubeService,
		subtitleService:   subtitleService,
		transcriptService: transcriptService,
		signalService:     signalService,
	}
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

	video, err := h.youtubeService.GetVideo((*c).Request().Context(), parsed.ID, metadata.ProviderYTDLP)
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

	windowSize := req.WindowSize
	if windowSize <= 0 {
		windowSize = 300
	}
	bucketCount := req.BucketCount
	if bucketCount <= 0 {
		if req.Analysis == "word-stats" {
			bucketCount = 32
		} else {
			bucketCount = 2
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

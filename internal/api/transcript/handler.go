package transcript

import (
	"net/http"

	"github.com/Arvind215271/askito/internal/api"
	youtubeurl "github.com/Arvind215271/askito/internal/youtube/input"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	youtubeService    *metadata.Service
	subtitleService   *subtitle.SubtitleService
	transcriptService *transcript.Service
}

func NewHandler(youtubeService *metadata.Service, subtitleService *subtitle.SubtitleService, transcriptService *transcript.Service) *Handler {
	return &Handler{
		youtubeService:    youtubeService,
		subtitleService:   subtitleService,
		transcriptService: transcriptService,
	}
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

	return (*c).JSON(http.StatusOK, t)
}

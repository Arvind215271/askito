package subtitle

import (
	"net/http"

	"github.com/Arvind215271/askito/internal/api"
	youtubeurl "github.com/Arvind215271/askito/internal/youtube/input"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	youtubeService  *metadata.Service
	subtitleService *subtitle.SubtitleService
}

func NewHandler(youtubeService *metadata.Service, subtitleService *subtitle.SubtitleService) *Handler {
	return &Handler{
		youtubeService:  youtubeService,
		subtitleService: subtitleService,
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

	video, err := h.youtubeService.GetVideo((*c).Request().Context(), parsed.ID, metadata.ProviderYTDLP)
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

	video, err := h.youtubeService.GetVideo((*c).Request().Context(), req.URL, metadata.ProviderYTDLP)
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

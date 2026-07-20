package video

import (
	"context"
	"net/http"

	"github.com/Arvind215271/askito/internal/api"
	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/description"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	youtubeurl "github.com/Arvind215271/askito/internal/youtube/input"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	youtubeService *metadata.Service
}

func NewHandler(youtubeService *metadata.Service) *Handler {
	return &Handler{
		youtubeService: youtubeService,
	}
}

func (h *Handler) GetVideoByID(c *echo.Context) error {
	var req VideoByIDRequest
	if err := (*c).Bind(&req); err != nil {
		return Err.BadRequest("Invalid request").Wrap(err)
	}

	if err := api.Validate(req); err != nil {
		return err
	}

	providerType := metadata.ProviderType(req.Provider)

	video, err := h.getVideo((*c).Request().Context(), req.ID, providerType)
	if err != nil {
		return Err.FetchFailed(err)
	}

	return (*c).JSON(http.StatusOK, VideoResponse{Video: video})
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

	providerType := metadata.ProviderType(req.Provider)

	video, err := h.getVideo((*c).Request().Context(), parsed.ID, providerType)
	if err != nil {
		return Err.FetchFailed(err)
	}

	return (*c).JSON(http.StatusOK, VideoResponse{Video: video})
}

func (h *Handler) getVideo(ctx context.Context, id string, providerType metadata.ProviderType) (youtube.Video, error) {
	if providerType == "" {
		providerType = metadata.ProviderYTDLP
	}

	if providerType != metadata.ProviderAPI && providerType != metadata.ProviderYTDLP {
		return youtube.Video{}, Err.InvalidProvider()
	}

	video, err := h.youtubeService.GetVideo(ctx, id, providerType)
	if err != nil {
		return youtube.Video{}, err
	}

	video.DescriptionMetadata = description.ProcessDescription(video.Description)
	return video, nil
}

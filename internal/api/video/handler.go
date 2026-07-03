package video

import (
	"context"
	"net/http"

	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/description"
	youtubeurl "github.com/Arvind215271/askito/internal/youtube/input"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	youtubeService *youtube.Service
}

func NewHandler(youtubeService *youtube.Service) *Handler {
	return &Handler{
		youtubeService: youtubeService,
	}
}

func (h *Handler) GetVideoByID(c *echo.Context) error {
	id := (*c).Param("id")
	if id == "" {
		return Err.IDRequired()
	}

	providerType := youtube.ProviderType((*c).QueryParam("provider"))
	
	video, err := h.getVideo( (*c).Request().Context(), id, providerType)
	if err != nil {
		return Err.FetchFailed(err)
	}

	return (*c).JSON(http.StatusOK, video)
}

func (h *Handler) GetVideoByURL(c *echo.Context) error {
	url := (*c).QueryParam("url")
	if url == "" {
		return Err.URLRequired()
	}

	parsed, err := youtubeurl.Parse(url)
	if err != nil {
		return Err.InvalidURL()
	}

	if parsed.InputType != youtubeurl.InputTypeVideo {
		return Err.NotAVideo()
	}

	providerType := youtube.ProviderType((*c).QueryParam("provider"))
	
	video, err := h.getVideo((*c).Request().Context(), parsed.ID, providerType)
	if err != nil {
		return Err.FetchFailed(err)
	}

	return (*c).JSON(http.StatusOK, video)
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

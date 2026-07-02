package video

import (
	"net/http"

	"github.com/Arvind215271/askito/internal/youtube"
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
		return (*c).JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	providerType := youtube.ProviderType((*c).QueryParam("provider"))
	if providerType == "" {
		providerType = youtube.ProviderYTDLP
	}

	if providerType != youtube.ProviderAPI && providerType != youtube.ProviderYTDLP {
		return (*c).JSON(http.StatusBadRequest, map[string]string{"error": "invalid provider"})
	}

	video, err := h.youtubeService.GetVideo((*c).Request().Context(), id, providerType)
	if err != nil {
		return (*c).JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return (*c).JSON(http.StatusOK, video)
}

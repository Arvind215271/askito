package playlist

import (
	"github.com/labstack/echo/v5"
)

func RegisterPlaylistRoutes(g *echo.Group, h *Handler) {
	g.POST("/videos", h.ExpandPlaylistVideos)
}

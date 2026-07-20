package subtitle

import (
	"github.com/labstack/echo/v5"
)

func RegisterSubtitleRoutes(g *echo.Group, h *Handler) {
	g.POST("/options", h.GetSubtitleOptions)
	g.POST("/download", h.DownloadSubtitle)
}

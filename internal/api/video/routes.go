package video

import (
	"github.com/labstack/echo/v5"
)

func RegisterVideoRoutes(g *echo.Group, h *Handler) {
	g.GET("/id", h.GetVideoByID)
	g.GET("/url", h.GetVideoByURL)
	g.POST("/transcripts", h.GetTranscript)
	g.POST("/signals", h.GetVideoSignals)
}

func RegisterSubtitleRoutes(g *echo.Group, h *Handler) {
	g.POST("/options", h.GetSubtitleOptions)
	g.POST("/download", h.DownloadSubtitle)
}

package export

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Group, handler *Handler) {
	e.POST("/video", handler.ExportVideo)
	e.POST("/videos", handler.ExportVideos)
	e.POST("/playlist", handler.ExportPlaylist)
}

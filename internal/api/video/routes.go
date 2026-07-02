package video

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, h *Handler) {
	g := e.Group("/videos")
	g.GET("/:id", h.GetVideoByID)
	g.GET("/url", h.GetVideoByURL)
}

package job

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Group, h *Handler) {
	e.GET("/:id", h.Get)
}

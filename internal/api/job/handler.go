package job

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/Arvind215271/askito/internal/api"
	domainJob "github.com/Arvind215271/askito/internal/job"
)

type Handler struct {
	manager *domainJob.JobManager
}

func NewHandler(manager *domainJob.JobManager) *Handler {
	return &Handler{
		manager: manager,
	}
}

func (h *Handler) GetUserJob(c *echo.Context) error {
	id := (*c).Param("id")

	job, err := h.manager.Get(id)
	if err != nil {
		if errors.Is(err, domainJob.ErrJobNotFound) {
			return Err.JobNotFound()
		}
		return api.Err.Common.Internal().Wrap(err)
	}

	return (*c).JSON(http.StatusOK, job)
}

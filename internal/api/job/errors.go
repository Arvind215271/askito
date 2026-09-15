package job

import (
	"net/http"

	"github.com/Arvind215271/askito/internal/api"
)

type JobErrors struct{}

var Err = JobErrors{}

func (JobErrors) JobNotFound() *api.AppError {
	return api.NewError(
		"JOB_NOT_FOUND",
		"Job not found",
		http.StatusNotFound,
	)
}

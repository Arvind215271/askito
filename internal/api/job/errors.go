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

func (JobErrors) MissingOrInvalidUserID() *api.AppError {
	return api.NewError(
		"INVALID_USER_ID",
		"Missing or invalid X-User-ID header",
		http.StatusBadRequest,
	)
}

func (JobErrors) ActiveJobExists() *api.AppError {
	return api.NewError(
		"ACTIVE_JOB_EXISTS",
		"User already has an active job",
		http.StatusConflict,
	)
}

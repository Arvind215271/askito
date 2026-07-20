package transcript

import (
	"net/http"

	"github.com/Arvind215271/askito/internal/api"
)

type TranscriptRequest struct {
	URL      string `json:"url" validate:"required,url"`
	Type     string `json:"type" validate:"required,oneof=manual automatic"`
	Language string `json:"language" validate:"required"`
}

type TranscriptErrors struct{}

var Err = TranscriptErrors{}

func (TranscriptErrors) BadRequest(msg string) *api.AppError {
	return api.NewError("BAD_REQUEST", msg, http.StatusBadRequest)
}

func (TranscriptErrors) InvalidURL() *api.AppError {
	return api.NewError("INVALID_URL", "Invalid URL", http.StatusBadRequest)
}

func (TranscriptErrors) FetchFailed(err error) *api.AppError {
	return api.NewError("FETCH_FAILED", "Failed to fetch video", http.StatusInternalServerError).Wrap(err)
}

func (TranscriptErrors) InternalError(err error) *api.AppError {
	return api.NewError("INTERNAL_ERROR", "Something went wrong", http.StatusInternalServerError).Wrap(err)
}

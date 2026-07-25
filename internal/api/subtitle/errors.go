package subtitle

import (
	"net/http"

	"github.com/Arvind215271/askito/internal/api"
)

type SubtitleErrors struct{}

var Err = SubtitleErrors{}

func (SubtitleErrors) BadRequest(msg string) *api.AppError {
	return api.NewError("BAD_REQUEST", msg, http.StatusBadRequest)
}

func (SubtitleErrors) InvalidURL() *api.AppError {
	return api.NewError("INVALID_URL", "Invalid URL", http.StatusBadRequest)
}

func (SubtitleErrors) FetchFailed(err error) *api.AppError {
	return api.NewError("FETCH_FAILED", "Failed to fetch video", http.StatusInternalServerError).Wrap(err)
}

func (SubtitleErrors) InvalidPreference(err error) *api.AppError {
	return api.NewError("INVALID_PREFERENCE", "Invalid subtitle preference", http.StatusBadRequest).Wrap(err)
}

func (SubtitleErrors) SubtitleNotFound(err error) *api.AppError {
	return api.NewError("SUBTITLE_NOT_FOUND", "No matching subtitle found", http.StatusNotFound).Wrap(err)
}

func (SubtitleErrors) InternalError(err error) *api.AppError {
	return api.NewError("INTERNAL_ERROR", "Something went wrong", http.StatusInternalServerError).Wrap(err)
}

package video

import (
	"net/http"

	"github.com/Arvind215271/askito/internal/api"
)

type VideoErrors struct{}

var Err = VideoErrors{}

func (VideoErrors) IDRequired() *api.AppError {
	return api.NewError(
		"VIDEO_ID_REQUIRED",
		"Video ID is required",
		http.StatusBadRequest,
	)
}

func (VideoErrors) URLRequired() *api.AppError {
	return api.NewError(
		"VIDEO_URL_REQUIRED",
		"Video URL is required",
		http.StatusBadRequest,
	)
}

func (VideoErrors) InvalidProvider() *api.AppError {
	return api.NewError(
		"INVALID_PROVIDER",
		"Invalid provider",
		http.StatusBadRequest,
	)
}

func (VideoErrors) InvalidURL() *api.AppError {
	return api.NewError(
		"INVALID_URL",
		"Invalid URL",
		http.StatusBadRequest,
	)
}

func (VideoErrors) NotAVideo() *api.AppError {
	return api.NewError(
		"NOT_A_VIDEO",
		"The provided URL is not a video",
		http.StatusBadRequest,
	)
}

func (VideoErrors) FetchFailed(err error) *api.AppError {
	return api.NewError(
		"FETCH_FAILED",
		"Failed to fetch video",
		http.StatusInternalServerError,
	).Wrap(err)
}

func (VideoErrors) BadRequest(msg string) *api.AppError {
	return api.NewError(
		"BAD_REQUEST",
		msg,
		http.StatusBadRequest,
	)
}

func (VideoErrors) InternalError(err error) *api.AppError {
	return api.NewError(
		"INTERNAL_ERROR",
		"Something went wrong",
		http.StatusInternalServerError,
	).Wrap(err)
}

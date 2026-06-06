// ./internal/youtube/errors.go
package youtube

import (
	"net/http"

	"github.com/Arvind215271/askito/internal/api"
)

type Errors struct {
	Playlist PlaylistErrors
	Video    VideoErrors
	Export   ExportErrors
}

var Err = Errors{
	Playlist: PlaylistErrors{},
	Video:    VideoErrors{},
	Export:   ExportErrors{},
}

type PlaylistErrors struct{}
type VideoErrors struct{}
type ExportErrors struct{}

// Playlist Errors

func (PlaylistErrors) InvalidURL() *api.AppError {
	return api.NewError(
		"PLAYLIST_INVALID_URL",
		"Invalid YouTube playlist URL",
		http.StatusBadRequest,
	)
}

func (PlaylistErrors) InvalidDomain() *api.AppError {
	return api.NewError(
		"PLAYLIST_INVALID_DOMAIN",
		"URL must be a YouTube URL",
		http.StatusBadRequest,
	)
}

func (PlaylistErrors) MissingID() *api.AppError {
	return api.NewError(
		"PLAYLIST_MISSING_ID",
		"Playlist ID is missing",
		http.StatusBadRequest,
	)
}

func (PlaylistErrors) InvalidType() *api.AppError {
	return api.NewError(
		"PLAYLIST_INVALID_TYPE",
		"URL is not a playlist",
		http.StatusBadRequest,
	)
}


func (PlaylistErrors) NotFound() *api.AppError {
	return api.NewError(
		"PLAYLIST_NOT_FOUND",
		"Playlist not found",
		http.StatusNotFound,
	)
}

func (PlaylistErrors) FetchFailed() *api.AppError {
	return api.NewError(
		"PLAYLIST_FETCH_FAILED",
		"Failed to fetch playlist",
		http.StatusInternalServerError,
	)
}

// video Errors

func (VideoErrors) NotFound() *api.AppError {
	return api.NewError(
		"VIDEO_NOT_FOUND",
		"Video not found",
		http.StatusNotFound,
	)
}

func (VideoErrors) FetchFailed() *api.AppError {
	return api.NewError(
		"VIDEO_FETCH_FAILED",
		"Failed to fetch video",
		http.StatusInternalServerError,
	)
}


// export errors
func (ExportErrors) InvalidFormat() *api.AppError {
	return api.NewError(
		"EXPORT_INVALID_FORMAT",
		"Invalid export format",
		http.StatusBadRequest,
	)
}

func (ExportErrors) MarshalFailed() *api.AppError {
	return api.NewError(
		"EXPORT_MARSHAL_FAILED",
		"Failed to export data",
		http.StatusInternalServerError,
	)
}

func (ExportErrors) InvalidField() *api.AppError {
	return api.NewError(
		"EXPORT_INVALID_FIELD",
		"One or more export fields are invalid",
		http.StatusBadRequest,
	)
}
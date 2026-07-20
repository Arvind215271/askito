package signal

import (
	"net/http"

	"github.com/Arvind215271/askito/internal/api"
	"github.com/Arvind215271/askito/internal/youtube/signal/word_stats"
)

type SignalRequest struct {
	URL        string  `json:"url" validate:"required"`
	Analysis   string  `json:"analysis" validate:"required,oneof=word-stats window-stats"`
	Type       string  `json:"type" validate:"required,oneof=manual automatic"`
	Language   string  `json:"language" validate:"required"`
	UseHeavy   bool    `json:"use_heavy_stopwords"`
	MinFreq    int     `json:"min_freq" validate:"gte=0"`
	Depth      float64 `json:"depth" validate:"gte=0,lte=1"`
	WindowSize float64 `json:"window_size" validate:"gt=0"`
	BucketCount int    `json:"bucket_count" validate:"gt=0"`
}

type SignalResponse struct {
	URL         string                   `json:"url"`
	WordStats   *wordstats.Result        `json:"word_stats,omitempty"`
	WindowStats []wordstats.Result       `json:"window_stats,omitempty"`
}

type SignalErrors struct{}

var Err = SignalErrors{}

func (SignalErrors) BadRequest(msg string) *api.AppError {
	return api.NewError("BAD_REQUEST", msg, http.StatusBadRequest)
}

func (SignalErrors) InvalidURL() *api.AppError {
	return api.NewError("INVALID_URL", "Invalid URL", http.StatusBadRequest)
}

func (SignalErrors) FetchFailed(err error) *api.AppError {
	return api.NewError("FETCH_FAILED", "Failed to fetch video", http.StatusInternalServerError).Wrap(err)
}

func (SignalErrors) InternalError(err error) *api.AppError {
	return api.NewError("INTERNAL_ERROR", "Something went wrong", http.StatusInternalServerError).Wrap(err)
}

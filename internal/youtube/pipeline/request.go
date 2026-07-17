package pipeline

import (
	"github.com/Arvind215271/askito/internal/youtube/signal"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
)

// Request composes feature-specific requests from each feature package.
// The pipeline only coordinates; validation and processing logic belongs to feature packages.
type Request struct {
	// Fields specifies which top-level fields to include in the response.
	// If empty, all fields are included.
	Fields []string

	// Subtitle request for downloading subtitles.
	// Owned by the subtitle package; validation happens there.
	Subtitle *subtitle.DownloadRequest

	// Transcript request for processing an already-parsed transcript.
	// Owned by the transcript package; validation happens there.
	Transcript *transcript.ProcessingRequest

	// Signal request for signal analysis on a transcript.
	// Owned by the signal package; validation happens there.
	Signal *signal.SignalRequest
}

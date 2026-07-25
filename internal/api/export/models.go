package export

import (
	subtitleapi "github.com/Arvind215271/askito/internal/api/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/signal"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
)

type CommonExportFields struct {
	Subtitle    *subtitle.DownloadRequest                 `json:"subtitle,omitempty"`
	Preferences []subtitleapi.SubtitlePreferenceRequest `json:"preferences,omitempty"`
	Transcript  *transcript.ProcessingRequest             `json:"transcript,omitempty"`
	Signal      *signal.SignalRequest                     `json:"signal,omitempty"`
	Format      string                                    `json:"format"`
	Fields      []string                                  `json:"fields,omitempty"`
}

type PlaylistExportRequest struct {
	CommonExportFields
	Input string `json:"input"`
}

type VideosExportRequest struct {
	CommonExportFields
	Inputs []string `json:"inputs"`
}

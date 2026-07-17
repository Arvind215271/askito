package export

import (
	"github.com/Arvind215271/askito/internal/youtube/signal"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
)

type VideoExportRequest struct {
	Input  string   `json:"input"`
	Fields []string `json:"fields,omitempty"`
	Format string   `json:"format"`

	Subtitle   *subtitle.DownloadRequest    `json:"subtitle,omitempty"`
	Transcript *transcript.ProcessingRequest `json:"transcript,omitempty"`
	Signal     *signal.SignalRequest       `json:"signal,omitempty"`
}

type PlaylistExportRequest struct {
	Input       string   `json:"input"`
	VideoFields []string `json:"video_fields,omitempty"`
	Format      string   `json:"format"`

	Subtitle   *subtitle.DownloadRequest      `json:"subtitle,omitempty"`
	Transcript *transcript.ProcessingRequest  `json:"transcript,omitempty"`
	Signal     *signal.SignalRequest         `json:"signal,omitempty"`
}
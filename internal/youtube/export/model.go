package export

import (
	"github.com/Arvind215271/askito/internal/youtube/fields"
)

type Format string

const (
	FormatJSON Format = "json"
)

type ExportData map[string]any

type PlaylistExportRequest struct {
	PlaylistID  string          `json:"playlist_id"`
	VideoFields *fields.Planner `json:"-"`
	Format      Format          `json:"format"`
}

type VideoExportRequest struct {
	VideoID string          `json:"video_id"`
	Fields  *fields.Planner `json:"-"`
	Format  Format          `json:"format"`
}

type BatchVideoExportRequest struct {
	VideoIDs    []string        `json:"video_ids"`
	VideoFields *fields.Planner `json:"-"`
	Format      Format          `json:"format"`
}

type ExportResponse struct {
	SchemaVersion string     `json:"schema_version"`
	Format        Format     `json:"format"`
	Data          ExportData `json:"data"`
}

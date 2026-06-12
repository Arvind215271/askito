package export



type Format string

const (
	FormatJSON Format = "json"
)

type ExportData map[string]any

type PlaylistExportRequest struct {
    PlaylistID string   `json:"playlist_id"`
    VideoFields []string `json:"video_fields,omitempty"`
    Format Format `json:"format"`
}

type VideoExportRequest struct {
	VideoID string   `json:"video_id"`
	Fields  []string `json:"fields,omitempty"`
	Format  Format   `json:"format"`
}

type ExportResponse struct {
	SchemaVersion string     `json:"schema_version"`
	Format        Format     `json:"format"`
	Data          ExportData `json:"data"`
}


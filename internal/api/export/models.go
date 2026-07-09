package export

type VideoExportRequest struct {
	Input  string   `json:"input"`
	Fields []string `json:"fields,omitempty"`
	Format string   `json:"format"`
}

type PlaylistExportRequest struct {
	Input       string   `json:"input"`
	VideoFields []string `json:"video_fields,omitempty"`
	Format      string   `json:"format"`
}

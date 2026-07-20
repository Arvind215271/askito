package subtitle

type SubtitleOptionsRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type SubtitleDownloadRequest struct {
	URL      string `json:"url" validate:"required,url"`
	Type     string `json:"type" validate:"required,oneof=manual automatic"`
	Language string `json:"language" validate:"required"`
	Format   string `json:"format" validate:"omitempty"`
}

package subtitle

type SubtitleOptionsRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type SubtitlePreferenceRequest struct {
	Language string `json:"language" validate:"required"`
	Type     string `json:"type" validate:"required,oneof=manual automatic manual>automatic automatic>manual"`
}

type SubtitleDownloadRequest struct {
	URL         string                      `json:"url" validate:"required,url"`
	Preferences []SubtitlePreferenceRequest `json:"preferences"`
	Format      string                      `json:"format,omitempty" validate:"omitempty,oneof=json3 vtt"`
}

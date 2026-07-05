package video

type SubtitleOptionsRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type SubtitleDownloadRequest struct {
	URL      string `json:"url" validate:"required,url"`
	Type     string `json:"type" validate:"required,oneof=manual automatic"`
	Language string `json:"language" validate:"required"`
	Format   string `json:"format" validate:"omitempty"`
}

type TranscriptRequest struct {
	URL      string `json:"url" validate:"required,url"`
	Type     string `json:"type" validate:"required,oneof=manual automatic"`
	Language string `json:"language" validate:"required"`
}

type VideoRequest struct {
	URL      string `json:"url" validate:"required,url"`
	Provider string `json:"provider" validate:"omitempty"`
}

type VideoByIDRequest struct {
	ID       string `json:"id" validate:"required"`
	Provider string `json:"provider" validate:"omitempty"`
}

type VideoResponse struct {
	Video interface{} `json:"video"`
}

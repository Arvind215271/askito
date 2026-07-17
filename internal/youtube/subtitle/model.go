package subtitle


type SubtitleTrack struct {
	LanguageCode string   `json:"languageCode"`
	LanguageName string   `json:"languageName"`
	Formats      []string `json:"formats"`
}

type SubtitleMetadata struct {
	Manual    []SubtitleTrack `json:"manual"`
	Automatic []SubtitleTrack `json:"automatic"`
}

type DownloadRequest struct {
	VideoID  string `json:"videoId"`
	Type     string `json:"type"` // manual | automatic
	Language string `json:"language"`
	Format   string `json:"format"`
}


type SubtitleResult struct {
	Content  []byte `json:"-"`
	Language string `json:"language"`
	Format   string `json:"format"`
}

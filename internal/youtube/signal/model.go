package signal

type SignalRequest struct {
	VideoID    string  `json:"video_id"`
	Analysis   string  `json:"analysis"` // "word-stats" or "window-stats"
	UseHeavy   bool    `json:"use_heavy_stopwords"`
	MinFreq    int     `json:"min_freq"`
	Depth      float64 `json:"depth"`
	WindowSize float64 `json:"window_size"`
}

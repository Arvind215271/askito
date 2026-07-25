package subtitle

func DefaultDownloadRequest(videoID string) DownloadRequest {
	return DownloadRequest{
		VideoID:  videoID,
		Language: "en",
		Type:     "automatic",
		Format:   "json3",
	}
}

func DefaultPreferences() []SubtitlePreference {
	return []SubtitlePreference{
		{
			Language: "en*",
			Type:     ManualPreferred,
		},
		{
			Language: "*",
			Type:     ManualPreferred,
		},
	}
}
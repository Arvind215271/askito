package subtitle

func DefaultDownloadRequest(videoID string) DownloadRequest {
	return DownloadRequest{
		VideoID:  videoID,
		Language: "en",
		Type:     "automatic",
		Format:   "json3",
	}
}
